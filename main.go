package main

import (
	"log"

	"github.com/hsmtkk/deltahedgesim/cmd"
)

func main() {
	cmd := cmd.RootCommand
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
