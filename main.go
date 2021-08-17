package main

import (
	"context"
	"fmt"

	"github.com/driver005/oauth/driver"
)

func main() {

	d := driver.New(context.Background())
	fmt.Println(d)
}
