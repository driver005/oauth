package main

import (
	"github.com/driver005/oauth/driver"
	"github.com/ory/x/logrusx"
)

func main() {
	l := logrusx.New("test", "master")
	driver.Init(l)
}
