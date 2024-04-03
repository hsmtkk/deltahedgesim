package hedge

import (
	"log"
	"os"

	"github.com/hsmtkk/aukabucomgo/base"
	"github.com/hsmtkk/deltahedgesim/hedge"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use: "hedge",
	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	apiPassword := os.Getenv("API_PASSWORD")
	if apiPassword == "" {
		log.Fatal("env var API_PASSWORD is not defined")
	}
	baseClient, err := base.New(base.PRODUCTION, apiPassword)
	if err != nil {
		log.Fatal(err)
	}
	if err := hedge.Hedge(baseClient); err != nil {
		log.Fatal(err)
	}
}
