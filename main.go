package main

import (
	"context"
	"fmt"

	"github.com/driver005/oauth/driver"
)

//go:generate swagger generate spec

func main() {
	newMigrateHandler().MigrateSQL()
	d := driver.New(context.Background(), driver.DisablePreloading())
	fmt.Println(d)
}
