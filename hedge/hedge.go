package hedge

import (
	"fmt"
	"slices"

	"github.com/hsmtkk/aukabucomgo/base"
	"github.com/hsmtkk/aukabucomgo/info/boardget"
	"github.com/hsmtkk/aukabucomgo/info/symbolnamefutureget"
	"github.com/hsmtkk/deltahedgesim/yaml/future"
	"github.com/hsmtkk/deltahedgesim/yaml/option"
	"github.com/hsmtkk/deltahedgesim/yaml/profitloss"
)

func Hedge(baseClient base.Client) error {
	totalDelta, err := CalcDisplayTotalDelta(baseClient)
	if err != nil {
		return err
	}

	direction := DecideHedge(totalDelta)
	switch direction {
	case NO_HEDGE:
		fmt.Println("No hedge is needed")
		return nil
	case BUY:
		fmt.Println("Buy back future")
		Buy(baseClient)
	case SELL:
		fmt.Println("Sell future")
		Sell(baseClient)
	}

	_, err = CalcDisplayTotalDelta(baseClient)
	if err != nil {
		return err
	}

	return nil
}

func OptionDelta(baseClient base.Client) (int, float64, float64, error) {
	optionPositions, err := option.Load()
	if err != nil {
		return 0, 0, 0, err
	}
	boardClient := boardget.New(baseClient)
	boardResp, err := boardClient.BoardGet(optionPositions.Symbol, boardget.ALL_DAY)
	if err != nil {
		return 0, 0, 0, err
	}
	singleDelta := boardResp.Delta
	quantity := optionPositions.Quantity
	optionDelta := singleDelta * float64(quantity)
	return quantity, singleDelta, optionDelta, nil
}

func FutureDelta(baseClient base.Client) (int, float64, error) {
	futurePositions, err := future.Load()
	if err != nil {
		return 0, 0, err
	}
	quantity := len(futurePositions.SoldPrices)
	futureDelta := float64(quantity) * -0.1 // マイクロ、売り
	return quantity, futureDelta, err
}

type HedgeDirection int

const (
	BUY HedgeDirection = iota
	SELL
	NO_HEDGE
)

func DecideHedge(totalDelta float64) HedgeDirection {
	if totalDelta > 0.1 {
		return SELL
	} else if totalDelta < -0.1 {
		return BUY
	}
	return NO_HEDGE
}

func Buy(baseClient base.Client) (int, error) {
	// オプションポジション取得
	optionPositions, err := option.Load()
	if err != nil {
		return 0, err
	}
	year := optionPositions.Year
	month := optionPositions.Month

	// symbol取得
	symbolClient := symbolnamefutureget.New(baseClient)
	symbolResp, err := symbolClient.SymbolNameFutureGet(symbolnamefutureget.NK225micro, year, month)
	if err != nil {
		return 0, err
	}

	// Bid取得
	boardClient := boardget.New(baseClient)
	boardResp, err := boardClient.BoardGet(symbolResp.Symbol, boardget.ALL_DAY)
	if err != nil {
		return 0, err
	}
	bidPrice := int(boardResp.BidPrice)

	// 保有先物ロード
	futurePositions, err := future.Load()
	if err != nil {
		return 0, err
	}

	// 最高値で売り建てたものを買い戻し
	soldPrices := futurePositions.SoldPrices
	slices.Sort(soldPrices)
	soldPrice := soldPrices[len(soldPrices)-1]
	futurePositions.SoldPrices = soldPrices[:len(soldPrices)-1]
	if err := futurePositions.Save(); err != nil {
		return 0, err
	}

	// 損益記録
	profitLoss := soldPrice - bidPrice
	profitLosses, err := profitloss.Load()
	if err != nil {
		return 0, err
	}
	profitLosses.ProfitLosses = append(profitLosses.ProfitLosses, profitLoss)
	if err := profitLosses.Save(); err != nil {
		return 0, err
	}

	return bidPrice, nil
}

func Sell(baseClient base.Client) (int, error) {
	// オプションポジション取得
	optionPositions, err := option.Load()
	if err != nil {
		return 0, err
	}
	year := optionPositions.Year
	month := optionPositions.Month

	// symbol取得
	symbolClient := symbolnamefutureget.New(baseClient)
	symbolResp, err := symbolClient.SymbolNameFutureGet(symbolnamefutureget.NK225micro, year, month)
	if err != nil {
		return 0, err
	}

	// Ask取得
	boardClient := boardget.New(baseClient)
	boardResp, err := boardClient.BoardGet(symbolResp.Symbol, boardget.ALL_DAY)
	if err != nil {
		return 0, err
	}
	askPrice := int(boardResp.AskPrice)

	// 保有先物ロード&セーブ
	futurePositions, err := future.Load()
	if err != nil {
		return 0, err
	}
	futurePositions.SoldPrices = append(futurePositions.SoldPrices, askPrice)
	if err := futurePositions.Save(); err != nil {
		return 0, err
	}

	return askPrice, nil
}

func CalcDisplayTotalDelta(baseClient base.Client) (float64, error) {
	quantity, singleDelta, optionDelta, err := OptionDelta(baseClient)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Option delta: %d * %f = %f\n", quantity, singleDelta, optionDelta)

	quantity, futureDelta, err := FutureDelta(baseClient)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Future delta: %d * -0.1 = %f\n", quantity, futureDelta)

	totalDelta := optionDelta + futureDelta
	fmt.Printf("Total delta: %f\n", totalDelta)
	return totalDelta, nil
}
