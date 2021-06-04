package tradecred

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mahendraHegde/tradecred-notifier/config"
	"github.com/mahendraHegde/tradecred-notifier/errors"
)

type TradeCred struct {
	config             *config.TradeCredConfig
	token              string
	tokenUpdateFailCtr int
}

func tokenRefresh(obj *TradeCred) {
	timer := time.NewTicker(time.Minute * 30)
	for range timer.C {
		_, err := obj.updateToken("", "")
		if err != nil {
			log.Println("Failed to update token :=", err)
			obj.tokenUpdateFailCtr++
			if obj.tokenUpdateFailCtr >= 3 {
				log.Println("cacelling the updateToken job", obj.tokenUpdateFailCtr)
				timer.Stop()
				return
			}
		}
		log.Println("Token updated Successfuly >>>>>>>")
	}
}

//NewTradeCred builds TradeCred
func NewTradeCred(config *config.TradeCredConfig) *TradeCred {
	obj := &TradeCred{config: config, token: ""}
	go tokenRefresh(obj)
	return obj
}

func (this *TradeCred) updateToken(email, password string) (string, error) {
	values := map[string]string{"refresh_token": this.config.RefreshToken, "grant_type": "refresh_token", "otp_verification_id": "null"}
	if email != "" && password != "" {
		values = map[string]string{"email": email, "password": password, "grant_type": "password", "otp_verification_id": "null"}
	}

	body, err := json.Marshal(values)

	if err != nil {
		return "", err
	}

	resp, err := http.Post(this.config.Base+"/oauth/token", "application/json",
		bytes.NewBuffer(body))

	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		var meta map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&meta)
		return "", errors.ApiError{
			Status:  resp.StatusCode,
			Message: "Error while connecting to TradeCred",
			Meta:    meta,
		}
	}

	res := map[string]string{}

	json.NewDecoder(resp.Body).Decode(&res)

	this.token = res["access_token"]
	this.config.RefreshToken = res["refresh_token"]
	return this.token, nil
}

func (this *TradeCred) GetDeals(page int, email, password string) ([]Deal, error) {
	this.updateToken(email, password)
	url := this.config.Base + "/deals?page=" + strconv.Itoa(page) + "&resource_path=deals&merge_lease_deal=true&type=deal"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+this.token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var meta map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&meta)
		return nil, errors.ApiError{
			Status:  resp.StatusCode,
			Message: "Error while connecting to TradeCred",
			Meta:    meta,
		}
	}
	var res struct {
		Data []Deal `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&res)

	return res.Data, nil
}

//TRumzVy4DLKTLitnbGlS-oK72-82QxdcirJ8gAhodYY
