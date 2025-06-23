package main

import (
	"os"

	"github.com/stateless-minds/kubo/cmd/ipfs/kubo"
)

func main() {
	os.Exit(kubo.Start(kubo.BuildDefaultEnv))
}
