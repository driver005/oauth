package helpers

import (
	"context"

	"github.com/ory/x/errorsx"

	"golang.org/x/crypto/bcrypt"
)

const defaultBCryptWorkFactor = 12

// BCrypt implements a BCrypt hasher.
type BCrypt struct {
	c config
}

type config interface {
	BCryptCost() int
}

// NewBCrypt returns a new BCrypt instance.
func NewBCrypt(c config) *BCrypt {
	return &BCrypt{
		c: c,
	}
}

func (b *BCrypt) Hash(ctx context.Context, data []byte) ([]byte, error) {
	cf := b.c.BCryptCost()
	if cf == 0 {
		cf = defaultBCryptWorkFactor
	}
	s, err := bcrypt.GenerateFromPassword(data, cf)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}
	return s, nil
}

func (b *BCrypt) Compare(ctx context.Context, hash, data []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, data); err != nil {
		return errorsx.WithStack(err)
	}
	return nil
}
