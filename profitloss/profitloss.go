package profitloss

import (
	"fmt"

	"github.com/hsmtkk/aukabucomgo/base"
	"github.com/hsmtkk/aukabucomgo/info/boardget"
	"github.com/hsmtkk/aukabucomgo/info/symbolnamefutureget"
	"github.com/hsmtkk/deltahedgesim/yaml/future"
	"github.com/hsmtkk/deltahedgesim/yaml/option"
)

func CalcDisplayTotalProfitLoss(baseClient base.Client) error {
	optionProfitLoss, err := OptionProfitLoss(baseClient)
	if err != nil {
		return err
	}
	fmt.Printf("Option profit loss: %d\n", optionProfitLoss)

	futureProfitLoss, err := FutureProfitLoss(baseClient)
	if err != nil {
		return err
	}
	fmt.Printf("Future profit loss: %d\n", futureProfitLoss)

	totalProfitLoss := optionProfitLoss + futureProfitLoss
	fmt.Printf("Total profit loss: %d\n", totalProfitLoss)
	return nil
}

func OptionProfitLoss(baseClient base.Client) (int, error) {
	// オプションポジション取得
	optionPositions, err := option.Load()
	if err != nil {
		return 0, err
	}

	// 売却なので、買い板にぶつける Ask取得
	boardClient := boardget.New(baseClient)
	boardResp, err := boardClient.BoardGet(optionPositions.Symbol, boardget.ALL_DAY)
	if err != nil {
		return 0, err
	}

	profitLoss := 100 * optionPositions.Quantity * (int(boardResp.AskPrice) - optionPositions.BoughtPrice)
	return profitLoss, nil
}

func FutureProfitLoss(baseClient base.Client) (int, error) {
	// オプションポジション取得
	optionPositions, err := option.Load()
	if err != nil {
		return 0, err
	}
	year := optionPositions.Year
	month := optionPositions.Month

	// 先物シンボル取得
	symbolClient := symbolnamefutureget.New(baseClient)
	symbolResp, err := symbolClient.SymbolNameFutureGet(symbolnamefutureget.NK225micro, year, month)
	if err != nil {
		return 0, err
	}

	// 買い戻しなので、売り板にぶつける Bid取得
	boardClient := boardget.New(baseClient)
	boardResp, err := boardClient.BoardGet(symbolResp.Symbol, boardget.ALL_DAY)
	if err != nil {
		return 0, err
	}
	bidPrice := int(boardResp.BidPrice)

	// 先物ポジション取得
	futurePositions, err := future.Load()
	if err != nil {
		return 0, err
	}

	profitLoss := 0
	for _, soldPrice := range futurePositions.SoldPrices {
		profitLoss += 10 * (soldPrice - bidPrice)
	}
	return profitLoss, nil
}
