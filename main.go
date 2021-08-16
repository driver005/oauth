package main

import (
	"github.com/driver005/oauth/driver"
	helper "github.com/driver005/oauth/helpers"
)

func main() {
	l := helper.New("test", "master")
	driver.Init(l)
}
