package main

import (
	"fmt"
)

type VersionCmd struct{}

func (c *VersionCmd) Run() error {
	fmt.Printf("@customrealms/cli (crx) v%s\n", version)
	return nil
}
