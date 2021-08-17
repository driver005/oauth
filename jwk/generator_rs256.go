/*
 * Copyright
 * Add later
 */

package jwk

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v2"
)

type RS256Generator struct {
	KeyLength int
}

func (g *RS256Generator) Generate(id, use string) (*jose.JSONWebKeySet, error) {
	if g.KeyLength < 4096 {
		g.KeyLength = 4096
	}

	key, err := rsa.GenerateKey(rand.Reader, g.KeyLength)
	if err != nil {
		return nil, errors.Errorf("Could not generate key because %s", err)
	} else if err = key.Validate(); err != nil {
		return nil, errors.Errorf("Validation failed because %s", err)
	}

	if id == "" {
		id = uuid.New()
	}

	// jose does not support this...
	key.Precomputed = rsa.PrecomputedValues{}
	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Algorithm:                   "RS256",
				Key:                         key,
				Use:                         use,
				KeyID:                       Ider("private", id),
				Certificates:                []*x509.Certificate{},
				CertificateThumbprintSHA1:   []uint8{},
				CertificateThumbprintSHA256: []uint8{},
			},
			{
				Algorithm:                   "RS256",
				Use:                         use,
				Key:                         &key.PublicKey,
				KeyID:                       Ider("public", id),
				Certificates:                []*x509.Certificate{},
				CertificateThumbprintSHA1:   []uint8{},
				CertificateThumbprintSHA256: []uint8{},
			},
		},
	}, nil
}
