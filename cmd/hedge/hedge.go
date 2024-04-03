package hedge

import (
	"log"
	"os"
	"time"

	"github.com/hsmtkk/aukabucomgo/base"
	"github.com/hsmtkk/deltahedgesim/hedge"
	"github.com/spf13/cobra"
)

const LOOP_INTERVAL_SECONDS = 600

var Command = &cobra.Command{
	Use: "hedge",
	Run: run,
}

var loop bool

func init() {
	Command.Flags().BoolVar(&loop, "loop", false, "loop")
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
	if loop {
		for {
			if err := hedge.Hedge(baseClient); err != nil {
				log.Fatal(err)
			}
			time.Sleep(LOOP_INTERVAL_SECONDS * time.Second)
		}

	} else {
		if err := hedge.Hedge(baseClient); err != nil {
			log.Fatal(err)
		}
	}
}
