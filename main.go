package main

import "github.com/ruets/hermit/cmd"

var Version = "dev"

func init() {
	cmd.Version = Version
}

func main() {
	cmd.Execute()
}
