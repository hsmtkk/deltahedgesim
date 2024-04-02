package getsymbol

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hsmtkk/aukabucomgo/base"
	"github.com/hsmtkk/aukabucomgo/info/symbolnameoptionget"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:  "getsymbol put/call yyyymm strice-price",
	Args: cobra.ExactArgs(3),
	Run:  run,
}

func run(cmd *cobra.Command, args []string) {
	putOrCallStr := args[0]
	month := args[1]
	strikePriceStr := args[2]

	apiPassword := os.Getenv("API_PASSWORD")
	if apiPassword == "" {
		log.Fatal("env var API_PASSWORD is not defined")
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil {
		log.Fatalf("failed to parse %s as int: %v", month, err)
	}

	var putOrCall symbolnameoptionget.PutOrCall
	switch putOrCallStr {
	case "put":
		putOrCall = symbolnameoptionget.PUT
	case "call":
		putOrCall = symbolnameoptionget.CALL
	}

	strikePrice, err := strconv.Atoi(strikePriceStr)
	if err != nil {
		log.Fatalf("failed to parse %s as int: %v", strikePriceStr, err)
	}

	symbol, symbolName, err := getSymbol(apiPassword, monthInt, putOrCall, strikePrice)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(symbol)
	fmt.Println(symbolName)
}

func getSymbol(apiPassword string, month int, putOrCall symbolnameoptionget.PutOrCall, strikePrice int) (string, string, error) {
	baseClient, err := base.New(base.PRODUCTION, apiPassword)
	if err != nil {
		return "", "", err
	}
	optionClient := symbolnameoptionget.New(baseClient)
	resp, err := optionClient.SymbolNameOptionGet(symbolnameoptionget.NK225miniop, month, putOrCall, strikePrice)
	if err != nil {
		return "", "", err
	}
	return resp.Symbol, resp.SymbolName, nil
}
