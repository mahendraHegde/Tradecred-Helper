package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahendraHegde/tradecred-notifier/config"
	"github.com/mahendraHegde/tradecred-notifier/tradecred"
)

type Controller struct {
	TradecredService *tradecred.TradeCred
	Config           *config.Configurations
}

func (this Controller) GetFilteredDeals(c *gin.Context) {
	pages := []int{1, 2, 3}
	input := c.MustGet("query").(GetDealsQS)
	body := c.MustGet("body").(Credentials)
	headers := c.MustGet("headers").(Headers)
	err := Authenticator(headers, this.Config.Server.ApiKey)
	if err != nil {
		c.JSONP(http.StatusUnauthorized, err)
		return
	}
	if input.Days == 0 {
		input.Days = 200
	}
	if input.MaxAmount == 0 {
		input.MaxAmount = 100000
	}
	if input.Rate == 0.0 {
		input.Rate = 11.55
	}
	filterd := []tradecred.Deal{}
	for _, p := range pages {
		deals, err := this.TradecredService.GetDeals(p, body.Email, body.Password)
		if err != nil {
			c.JSONP(http.StatusInternalServerError, err)
			return
		}
		body.Email = ""
		body.Password = ""
		for _, deal := range deals {
			if deal.Attributes.State == "in_progress" && deal.Attributes.Days < float64(input.Days) && deal.Attributes.MinAmount < input.MaxAmount && deal.Attributes.MinAmount > 0 && deal.Attributes.Rate > input.Rate {
				filterd = append(filterd, deal)
			}
		}
	}

	c.JSON(http.StatusOK, filterd)
}
