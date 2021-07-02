package callmebot

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mahendraHegde/tradecred-notifier/config"
	"github.com/mahendraHegde/tradecred-notifier/errors"
	"github.com/mahendraHegde/tradecred-notifier/tradecred"
)

type CallMeBot struct {
	config    *config.CallMeBot
	sentDeals map[string]bool //dont send same message twice
}

//NewCallmeBot builds CallMeBot
func NewCallmeBot(config *config.CallMeBot) *CallMeBot {
	obj := &CallMeBot{config, map[string]bool{}}
	return obj
}

func (this *CallMeBot) SendWhatsAppMessage(ctx context.Context, deals []tradecred.Deal) error {
	text := "New Deals"
	dealsTobeSent := 0
	for i, deal := range deals {
		if _, ok := this.sentDeals[deal.Attributes.Code]; ok {
			log.Println("skipping sending deal to whats app as its already sent :", deal)
			continue
		}

		dealsTobeSent++

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
	if dealsTobeSent <= 0 {
		return nil
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
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to parse callmeBot response => ", err)
		return nil
	}
	stringBody := string(bodyBytes)
	if !strings.Contains(strings.ToLower(stringBody), "message queued") {
		log.Println("Callmebot error response " + stringBody)
	} else {
		for _, deal := range deals {
			this.sentDeals[deal.Attributes.Code] = true
		}

	}
	return nil
}
