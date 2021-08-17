/*
 * Copyright
 * Add later
 */

package jwk

import (
	"crypto/rand"
	"crypto/x509"
	"io"

	"github.com/ory/x/errorsx"

	"github.com/pborman/uuid"
	jose "gopkg.in/square/go-jose.v2"
)

type HS512Generator struct{}

func (g *HS512Generator) Generate(id, use string) (*jose.JSONWebKeySet, error) {
	// Taken from NewHMACKey
	key := &[32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	if id == "" {
		id = uuid.New()
	}

	var sliceKey = key[:]

	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Algorithm:    "HS512",
				Key:          sliceKey,
				Use:          use,
				KeyID:        id,
				Certificates: []*x509.Certificate{},
			},
		},
	}, nil
}
