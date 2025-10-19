package main

import (
	"fmt"

	"github.com/racg0092/rhombifer"
)

// Version of CLI tool
var VersionCmd = &rhombifer.Command{
	Name:      "version",
	ShortDesc: "application version",
	Run: func(args ...string) error {
		fmt.Println("v0.1.0")
		return nil
	},
}

func init() {
	rhombifer.Root().AddSub(VersionCmd)
}
