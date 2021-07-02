package core

import (
	"context"
	"log"

	"github.com/mahendraHegde/tradecred-notifier/tradecred"
)

func GetFilteredDeals(input GetDealsQS, body Credentials, TradecredService *tradecred.TradeCred, pages int) ([]tradecred.Deal, error) {
	if input.Days == 0 {
		input.Days = 200
	}
	if input.MaxAmount == 0 {
		input.MaxAmount = 170000
	}
	if input.Rate == 0.0 {
		input.Rate = 10.00
	}
	filterd := []tradecred.Deal{}
	_, err := TradecredService.UpdateToken(body.Email, body.Password)
	if err != nil {
		return nil, err
	}

	type ResErr struct {
		Deals []tradecred.Deal
		Err   error
	}
	ch := make(chan ResErr)
	defer close(ch)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processReponse := func(ctx context.Context, deals []tradecred.Deal, err error) {
		select {
		case <-ctx.Done():
			log.Println("[GetFilteredDeals]cancelling Get deals :=", ctx.Err())
			return
		default:
			log.Println("[GetFilteredDeals]checking response for filter ", input)
		}
		if err != nil {
			ch <- ResErr{Deals: nil, Err: err}
			return
		}
		rs := []tradecred.Deal{}
		for _, deal := range deals {
			if deal.Attributes.State == "in_progress" && deal.Attributes.Days < float64(input.Days) && deal.Attributes.MinAmount < input.MaxAmount && deal.Attributes.MinAmount > 0 && deal.Attributes.Rate > input.Rate {
				rs = append(rs, deal)
			}
		}
		ch <- ResErr{Deals: rs, Err: nil}
	}

	for p := 1; p <= pages; p++ {
		go func(ctx context.Context, p int) {
			deals, err := TradecredService.GetDeals(ctx, p)
			processReponse(ctx, deals, err)
		}(ctx, p)
	}

	go func(ctx context.Context) {
		deals, err := TradecredService.GetLiquidationRequests(ctx)
		processReponse(ctx, deals, err)
	}(ctx)

	for p := 0; p <= pages; p++ {
		rs := <-ch
		if rs.Err != nil {
			return nil, rs.Err
		}
		filterd = append(filterd, rs.Deals...)
	}
	return filterd, nil
}
