/*
 * Copyright
 * Add later
 */

package jwk

import jose "gopkg.in/square/go-jose.v2"

type KeyGenerator interface {
	Generate(id, use string) (*jose.JSONWebKeySet, error)
}
