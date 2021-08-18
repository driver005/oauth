package server

import (
	"context"
	"crypto/sha1" // #nosec G505 - This is required for certificate chains alongside sha256
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"

	"gopkg.in/square/go-jose.v2"

	"github.com/driver005/oauth/config"
	"github.com/driver005/oauth/registry"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ory/x/tlsx"

	"github.com/driver005/oauth/jwk"
)

const (
	tlsKeyName = "hydra.https-tls"
)

func AttachCertificate(priv *jose.JSONWebKey, cert *x509.Certificate) {
	priv.Certificates = []*x509.Certificate{cert}
	sig256 := sha256.Sum256(cert.Raw)
	// #nosec G401 - This is required for certificate chains alongside sha256
	sig1 := sha1.Sum(cert.Raw)
	priv.CertificateThumbprintSHA256 = sig256[:]
	priv.CertificateThumbprintSHA1 = sig1[:]
}

func GetOrCreateTLSCertificate(cmd *cobra.Command, d registry.Registry, iface config.ServeInterface) []tls.Certificate {
	cert, err := d.Config().TLS(iface).Certificate()

	if err == nil {
		return cert
	} else if !errors.Is(err, tlsx.ErrNoCertificatesConfigured) {
		d.Logger().WithError(err).Fatalf("Unable to load HTTPS TLS Certificate")
	}

	_, priv, err := jwk.AsymmetricKeypair(cmd.Context(), d, &jwk.RS256Generator{KeyLength: 4069}, tlsKeyName)
	if err != nil {
		d.Logger().WithError(err).Fatal("Unable to fetch HTTPS TLS key pairs")
	}

	if len(priv.Certificates) == 0 {
		cert, err := tlsx.CreateSelfSignedCertificate(priv.Key)
		if err != nil {
			d.Logger().WithError(err).Fatalf(`Could not generate a self signed TLS certificate`)
		}

		AttachCertificate(priv, cert)
		if err := d.KeyManager().DeleteKey(context.TODO(), tlsKeyName, priv.KeyID); err != nil {
			d.Logger().WithError(err).Fatal(`Could not update (delete) the self signed TLS certificate`)
		}

		if err := d.KeyManager().AddKey(context.TODO(), tlsKeyName, priv); err != nil {
			d.Logger().WithError(err).Fatalf(`Could not update (add) the self signed TLS certificate: %s %x %d`, cert.SignatureAlgorithm, cert.Signature, len(cert.Signature))
		}
	}

	block, err := jwk.PEMBlockForKey(priv.Key)
	if err != nil {
		d.Logger().WithError(err).Fatalf("Could not encode key to PEM")
	}

	if len(priv.Certificates) == 0 {
		d.Logger().Fatal("TLS certificate chain can not be empty")
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: priv.Certificates[0].Raw})
	pemKey := pem.EncodeToMemory(block)
	ct, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		d.Logger().WithError(err).Fatalf("Could not decode certificate")
	}

	return []tls.Certificate{ct}
}
