package main

import (
	"evaframe/cmd"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "evaframe"
	// Version is the version of the compiled software.
	Version string = "dev"
)

func main() {
	cmd.Execute()
}
