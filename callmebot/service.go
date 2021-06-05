package callmebot

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mahendraHegde/tradecred-notifier/config"
	"github.com/mahendraHegde/tradecred-notifier/errors"
	"github.com/mahendraHegde/tradecred-notifier/tradecred"
)

type CallMeBot struct {
	config *config.CallMeBot
}

//NewCallmeBot builds CallMeBot
func NewCallmeBot(config *config.CallMeBot) *CallMeBot {
	obj := &CallMeBot{config}
	return obj
}

func (this *CallMeBot) SendWhatsAppMessage(ctx context.Context, deals []tradecred.Deal) error {
	text := "New Deals"
	for i, deal := range deals {
		t := `
` + strconv.Itoa(i+1) + ". " + `
` + "code : " + deal.Attributes.Code

		if deal.Attributes.Name != "" {
			temp := `
Name : ` + deal.Attributes.Name
			t += temp
		}

		text += t
	}
	text = url.QueryEscape(text)
	url := this.config.WhatsApp.Base + "?phone=" + this.config.WhatsApp.Phone + "&text=" + text + "&apikey=" + this.config.WhatsApp.ApiKey
	log.Println("Message TEXT => " + text)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var meta map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&meta)
		return errors.ApiError{
			Status:  resp.StatusCode,
			Message: "Error while connecting to Callmebot",
			Meta:    meta,
		}
	}
	return nil
}
