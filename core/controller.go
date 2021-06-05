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
	input := c.MustGet("query").(GetDealsQS)
	body := c.MustGet("body").(Credentials)
	headers := c.MustGet("headers").(Headers)
	err := Authenticator(headers, this.Config.Server.ApiKey)
	if err != nil {
		c.JSONP(http.StatusUnauthorized, err)
		return
	}
	filtered, err := GetFilteredDeals(input, body, this.TradecredService, 3)
	if err != nil {
		c.JSONP(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, filtered)
}
