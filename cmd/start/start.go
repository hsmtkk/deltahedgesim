package start

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hsmtkk/aukabucomgo/base"
	"github.com/hsmtkk/aukabucomgo/info/boardget"
	"github.com/hsmtkk/aukabucomgo/info/symbolget"
	"github.com/hsmtkk/deltahedgesim/yaml/future"
	"github.com/hsmtkk/deltahedgesim/yaml/option"
	"github.com/hsmtkk/deltahedgesim/yaml/profitloss"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var Command = &cobra.Command{
	Use:  "start symbol quantity",
	Run:  run,
	Args: cobra.ExactArgs(2),
}

func run(cmd *cobra.Command, args []string) {
	symbol := args[0]
	quantityStr := args[1]
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		log.Fatalf("failed to parse %s as int: %v", quantityStr, err)
	}

	apiPassword := os.Getenv("API_PASSWORD")
	if apiPassword == "" {
		log.Fatal("env var API_PASSWORD is not defined")
	}

	baseClient, err := base.New(base.PRODUCTION, apiPassword)
	if err != nil {
		log.Fatal(err)
	}

	if err := initOption(baseClient, symbol, quantity); err != nil {
		log.Fatal(err)
	}
	if err := initFuture(); err != nil {
		log.Fatal(err)
	}
	if err := initProfitLoss(); err != nil {
		log.Fatal(err)
	}
}

func initOption(baseClient base.Client, symbol string, quantity int) error {
	boardClient := boardget.New(baseClient)
	boardResp, err := boardClient.BoardGet(symbol, boardget.ALL_DAY)
	if err != nil {
		return err
	}
	symbolClient := symbolget.New(baseClient)
	symbolResp, err := symbolClient.SymbolGet(symbol, symbolget.ALL_DAY)
	if err != nil {
		return err
	}
	date, err := time.Parse("2006/01", symbolResp.DerivMonth)
	if err != nil {
		return fmt.Errorf("failed to parse %s as yyyy/mm: %w", symbolResp.DerivMonth, err)
	}
	data := option.Schema{
		Symbol:      symbol,
		SymbolName:  boardResp.SymbolName,
		Year:        date.Year(),
		Month:       int(date.Month()),
		Quantity:    quantity,
		BoughtPrice: int(boardResp.BidPrice),
	}
	if err := writeYAML(data, "data/option.yaml"); err != nil {
		return err
	}
	return nil
}

func initFuture() error {
	data := future.Schema{
		SoldPrices: []int{},
	}
	if err := writeYAML(data, "data/future.yaml"); err != nil {
		return err
	}
	return nil
}

func initProfitLoss() error {
	data := profitloss.Schema{
		ProfitLosses: []int{},
	}
	if err := writeYAML(data, "data/profitloss.yaml"); err != nil {
		return err
	}
	return nil
}

func writeYAML(data interface{}, path string) error {
	bs, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	if err := os.WriteFile(path, bs, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
