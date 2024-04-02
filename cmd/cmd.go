package cmd

import (
	"github.com/hsmtkk/deltahedgesim/cmd/getsymbol"
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use: "deltahedgesim",
}

func init() {
	RootCommand.AddCommand(getsymbol.Command)
}
