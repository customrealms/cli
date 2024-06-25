package main

import (
	"fmt"

	"github.com/customrealms/cli/pkg/version"
)

type VersionCmd struct{}

func (c *VersionCmd) Run() error {
	fmt.Printf("@customrealms/cli (crx) v%s\n", version.Version)
	return nil
}
