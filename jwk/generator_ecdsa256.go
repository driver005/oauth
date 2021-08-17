/*
 * Copyright
 * Add later
 */

package jwk

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v2"
)

type ECDSA256Generator struct{}

func (g *ECDSA256Generator) Generate(id, use string) (*jose.JSONWebKeySet, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, errors.Errorf("Could not generate key because %s", err)
	}

	if id == "" {
		id = uuid.New()
	}

	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Key:          key,
				Use:          use,
				KeyID:        Ider("private", id),
				Certificates: []*x509.Certificate{},
			},
			{
				Key:          &key.PublicKey,
				Use:          use,
				KeyID:        Ider("public", id),
				Certificates: []*x509.Certificate{},
			},
		},
	}, nil
}
