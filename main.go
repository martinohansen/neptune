package main

import (
	"log"

	"github.com/martinohansen/neptune/cmd"
)

func init() {
	// Include file path for logging statements
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {
	cmd.Execute()
}
