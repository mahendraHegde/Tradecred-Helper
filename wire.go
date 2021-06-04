package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mahendraHegde/tradecred-notifier/config"
	"github.com/mahendraHegde/tradecred-notifier/core"
	"github.com/mahendraHegde/tradecred-notifier/tradecred"
)

var ctx = context.TODO()

//Wire wires all dependecies together
type Wire struct {
	config           *config.Configurations
	router           *gin.Engine
	tradeCredService *tradecred.TradeCred
	coreController   *core.Controller
}

func buildDependencies(config config.Configurations, router *gin.Engine) (*Wire, error) {
	tradeCredService := tradecred.NewTradeCred(&config.TradeCredConfig)
	coreController := core.Controller{TradecredService: tradeCredService, Config: &config}
	return &Wire{&config, router, tradeCredService, &coreController}, nil
}
