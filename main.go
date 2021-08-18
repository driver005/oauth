package main

import (
	"context"
	"fmt"

	"github.com/driver005/oauth/driver"
)

//go:generate swagger generate spec

func main() {

	d := driver.New(context.Background())
	fmt.Println(d)
}
