package main

import "github.com/driver005/oauth/cmd"

//go:generate swagger generate spec

func main() {
	cmd.Execute()
}
