package cmd

import (
	"github.com/hsmtkk/deltahedgesim/cmd/getsymbol"
	"github.com/hsmtkk/deltahedgesim/cmd/hedge"
	"github.com/hsmtkk/deltahedgesim/cmd/start"
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use: "deltahedgesim",
}

func init() {
	RootCommand.AddCommand(getsymbol.Command)
	RootCommand.AddCommand(start.Command)
	RootCommand.AddCommand(hedge.Command)
}
