package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ory/x/configx"

	analytics "github.com/ory/analytics-go/v4"

	"github.com/ory/x/reqlog"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/ory/x/healthx"
	"github.com/ory/x/metricsx"
	"github.com/ory/x/networkx"

	"github.com/driver005/oauth/client"
	"github.com/driver005/oauth/config"
	"github.com/driver005/oauth/consent"
	"github.com/driver005/oauth/driver"
	"github.com/driver005/oauth/helpers"
	"github.com/driver005/oauth/jwk"
	"github.com/driver005/oauth/oauth2"
	"github.com/driver005/oauth/registry"
	prometheus "github.com/ory/x/prometheusx"
)

var _ = &consent.Handler{}

func EnhanceMiddleware(d registry.Registry, n *negroni.Negroni, address string, router *httprouter.Router, enableCORS bool, iface config.ServeInterface) http.Handler {
	if !networkx.AddressIsUnixSocket(address) {
		n.UseFunc(helpers.RejectInsecureRequests(d, d.Config().TLS(iface)))
	}
	n.UseHandler(router)

	if !enableCORS {
		return n
	}

	options, enabled := d.Config().CORS(iface)
	if !enabled {
		return n
	}

	if enabled {
		d.Logger().
			WithField("options", fmt.Sprintf("%+v", options)).
			Infof("Enabling CORS on interface: %s", address)
		return cors.New(options).Handler(n)
	}
	return n
}

func isDSNAllowed(r registry.Registry) {
	if r.Config().DSN() == "memory" {
		r.Logger().Fatalf(`When using "hydra serve admin" or "hydra serve public" the DSN can not be set to "memory".`)
	}
}

func RunServeAdmin(cmd *cobra.Command, _ []string) {
	d := driver.New(cmd.Context(), driver.WithOptions(configx.WithFlags(cmd.Flags())))
	isDSNAllowed(d)

	admin, _, adminmw, _ := setup(d, cmd)
	cert := GetOrCreateTLSCertificate(cmd, d, config.AdminInterface) // we do not want to run this concurrently.

	d.PrometheusManager().RegisterRouter(admin.Router)

	var wg sync.WaitGroup
	wg.Add(1)

	go serve(
		d,
		cmd,
		&wg,
		config.AdminInterface,
		EnhanceMiddleware(d, adminmw, d.Config().ListenOn(config.AdminInterface), admin.Router, true, config.AdminInterface),
		d.Config().ListenOn(config.AdminInterface),
		d.Config().SocketPermission(config.AdminInterface),
		cert,
	)

	wg.Wait()
}

func RunServePublic(cmd *cobra.Command, _ []string) {
	d := driver.New(cmd.Context(), driver.WithOptions(configx.WithFlags(cmd.Flags())))
	isDSNAllowed(d)

	_, public, _, publicmw := setup(d, cmd)
	cert := GetOrCreateTLSCertificate(cmd, d, config.PublicInterface) // we do not want to run this concurrently.

	d.PrometheusManager().RegisterRouter(public.Router)

	var wg sync.WaitGroup
	wg.Add(1)

	go serve(
		d,
		cmd,
		&wg,
		config.PublicInterface,
		EnhanceMiddleware(d, publicmw, d.Config().ListenOn(config.PublicInterface), public.Router, false, config.PublicInterface),
		d.Config().ListenOn(config.PublicInterface),
		d.Config().SocketPermission(config.PublicInterface),
		cert,
	)

	wg.Wait()
}

func RunServeAll(cmd *cobra.Command, _ []string) {
	d := driver.New(cmd.Context(), driver.WithOptions(configx.WithFlags(cmd.Flags())), driver.DisablePreloading())

	fmt.Println(d)
	admin, public, adminmw, publicmw := setup(d, cmd)

	d.PrometheusManager().RegisterRouter(admin.Router)
	d.PrometheusManager().RegisterRouter(public.Router)

	var wg sync.WaitGroup
	wg.Add(2)

	GetOrCreateTLSCertificate(cmd, d, config.AdminInterface)

	go serve(
		d,
		cmd,
		&wg,
		config.PublicInterface,
		EnhanceMiddleware(d, publicmw, d.Config().ListenOn(config.PublicInterface), public.Router, false, config.PublicInterface),
		d.Config().ListenOn(config.PublicInterface),
		d.Config().SocketPermission(config.PublicInterface),
		GetOrCreateTLSCertificate(cmd, d, config.PublicInterface),
	)

	go serve(
		d,
		cmd,
		&wg,
		config.AdminInterface,
		EnhanceMiddleware(d, adminmw, d.Config().ListenOn(config.AdminInterface), admin.Router, true, config.AdminInterface),
		d.Config().ListenOn(config.AdminInterface),
		d.Config().SocketPermission(config.AdminInterface),
		GetOrCreateTLSCertificate(cmd, d, config.AdminInterface),
	)

	wg.Wait()
}

func setup(d registry.Registry, cmd *cobra.Command) (admin *helpers.RouterAdmin, public *helpers.RouterPublic, adminmw, publicmw *negroni.Negroni) {
	fmt.Println(banner(config.Version))

	if d.Config().CGroupsV1AutoMaxProcsEnabled() {
		_, err := maxprocs.Set(maxprocs.Logger(d.Logger().Infof))

		if err != nil {
			d.Logger().WithError(err).Fatal("Couldn't set GOMAXPROCS")
		}
	}

	adminmw = negroni.New()
	publicmw = negroni.New()

	admin = helpers.NewRouterAdmin()
	public = helpers.NewRouterPublic()

	if tracer := d.Tracer(cmd.Context()); tracer.IsLoaded() {
		adminmw.Use(tracer)
		publicmw.Use(tracer)
	}

	adminLogger := reqlog.
		NewMiddlewareFromLogger(d.Logger(),
			fmt.Sprintf("hydra/admin: %s", d.Config().IssuerURL().String()))
	if d.Config().DisableHealthAccessLog(config.AdminInterface) {
		adminLogger = adminLogger.ExcludePaths(healthx.AliveCheckPath, healthx.ReadyCheckPath)
	}

	adminmw.Use(adminLogger)
	adminmw.Use(d.PrometheusManager())

	publicLogger := reqlog.NewMiddlewareFromLogger(
		d.Logger(),
		fmt.Sprintf("hydra/public: %s", d.Config().IssuerURL().String()),
	)
	if d.Config().DisableHealthAccessLog(config.PublicInterface) {
		publicLogger.ExcludePaths(healthx.AliveCheckPath, healthx.ReadyCheckPath)
	}

	publicmw.Use(publicLogger)
	publicmw.Use(d.PrometheusManager())

	metrics := metricsx.New(
		cmd,
		d.Logger(),
		d.Config().Source(),
		&metricsx.Options{
			Service: "ory-hydra",
			ClusterID: metricsx.Hash(fmt.Sprintf("%s|%s",
				d.Config().IssuerURL().String(),
				d.Config().DSN(),
			)),
			IsDevelopment: d.Config().DSN() == "memory" ||
				d.Config().IssuerURL().String() == "" ||
				strings.Contains(d.Config().IssuerURL().String(), "localhost"),
			WriteKey: "h8dRH3kVCWKkIFWydBmWsyYHR4M0u0vr",
			WhitelistedPaths: []string{
				jwk.KeyHandlerPath,
				jwk.WellKnownKeysPath,

				client.ClientsHandlerPath,

				oauth2.DefaultConsentPath,
				oauth2.DefaultLoginPath,
				oauth2.DefaultPostLogoutPath,
				oauth2.DefaultLogoutPath,
				oauth2.DefaultErrorPath,
				oauth2.TokenPath,
				oauth2.AuthPath,
				oauth2.LogoutPath,
				oauth2.UserinfoPath,
				oauth2.WellKnownPath,
				oauth2.JWKPath,
				oauth2.IntrospectPath,
				oauth2.RevocationPath,
				oauth2.FlushPath,

				consent.ConsentPath,
				consent.ConsentPath + "/accept",
				consent.ConsentPath + "/reject",
				consent.LoginPath,
				consent.LoginPath + "/accept",
				consent.LoginPath + "/reject",
				consent.LogoutPath,
				consent.LogoutPath + "/accept",
				consent.LogoutPath + "/reject",
				consent.SessionsPath + "/login",
				consent.SessionsPath + "/consent",

				healthx.AliveCheckPath,
				healthx.ReadyCheckPath,
				healthx.VersionPath,
				prometheus.MetricsPrometheusPath,
				"/",
			},
			BuildVersion: config.Version,
			BuildTime:    config.Date,
			BuildHash:    config.Commit,
			Config: &analytics.Config{
				Endpoint:             "https://sqa.ory.sh",
				GzipCompressionLevel: 6,
				BatchMaxSize:         500 * 1000,
				BatchSize:            250,
				Interval:             time.Hour * 24,
			},
		},
	)

	adminmw.Use(metrics)
	publicmw.Use(metrics)

	d.RegisterRoutes(admin, public)

	return
}

func certsetup() (serverTLSConf *tls.Config, clientTLSConf *tls.Config, err error) {
	// set up our CA certificate
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(caPEM.Bytes())
	clientTLSConf = &tls.Config{
		RootCAs: certpool,
	}

	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		return nil, nil, err
	}

	serverTLSConf = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	return serverTLSConf, clientTLSConf, nil
}

func serve(
	d registry.Registry,
	cmd *cobra.Command,
	wg *sync.WaitGroup,
	iface config.ServeInterface,
	handler http.Handler,
	address string,
	permission *configx.UnixPermission,
	cert []tls.Certificate,
) {
	defer wg.Done()

	s := &http.Server{
		Handler:        handler,
		TLSConfig:      &tls.Config{Certificates: cert},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if d.Tracer(cmd.Context()).IsLoaded() {
		s.RegisterOnShutdown(d.Tracer(cmd.Context()).Close)
	}

	d.Logger().Infof("Setting up http server on %s", address)
	listener, err := networkx.MakeListener(address, permission)
	if err != nil {
		d.Logger().WithError(err).Fatal("Couldn't not create Listener")
	}

	s.ServeTLS(listener, "", "")

	// s.ListenAndServeTLS("", "")
}
