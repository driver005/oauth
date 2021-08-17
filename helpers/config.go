package helpers

import "net/url"

func UrlRoot(u *url.URL) *url.URL {
	if u.Path == "" {
		u.Path = "/"
	}
	return u
}
