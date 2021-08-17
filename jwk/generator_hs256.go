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

type HS256Generator struct{}

func (g *HS256Generator) Generate(id, use string) (*jose.JSONWebKeySet, error) {
	// Taken from NewHMACKey
	key := &[16]byte{}
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
				Algorithm:    "HS256",
				Use:          use,
				Key:          sliceKey,
				KeyID:        id,
				Certificates: []*x509.Certificate{},
			},
		},
	}, nil
}
