package job

import (
	"context"
	"log"
	"time"

	"github.com/mahendraHegde/tradecred-notifier/callmebot"
	"github.com/mahendraHegde/tradecred-notifier/core"
	"github.com/mahendraHegde/tradecred-notifier/tradecred"
)

func getCurrentHour() int {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Println(`Failed to load time`, err)
	}
	return time.Now().In(loc).Hour()
}

func ScheduleDealsCheck(ctx context.Context, duration time.Duration, TradecredService *tradecred.TradeCred, CallmeBotService *callmebot.CallMeBot) {
	log.Println("started ScheduleDealsCheck...")
	timer := time.NewTicker(duration)
	for range timer.C {
		select {
		case <-ctx.Done():
			log.Println("stopping ScheduleDealsCheck :=", ctx.Err())
			timer.Stop()
			return
		default:
			log.Println("starting ScheduleDealsCheck...")
		}
		currentHour := getCurrentHour()
		pages := 3
		if currentHour > 1 && currentHour < 6 {
			pages = 0
			log.Println("skipping ScheduleDealsCheck as its DND hour")
		}
		res, err := core.GetFilteredDeals(core.GetDealsQS{}, core.Credentials{}, TradecredService, pages)
		if err != nil {
			log.Println("Failed to run ScheduleDealsCheck :=", err)
		} else {
			log.Println("[ScheduleDealsCheck]Deals found >>>>>>>", res)
			if len(res) > 0 {
				err := CallmeBotService.SendWhatsAppMessage(ctx, res)
				if err != nil {
					log.Println("[ScheduleDealsCheck]Failed to send Whatsapp notification:=", err)
				}
			}
		}
	}
}
