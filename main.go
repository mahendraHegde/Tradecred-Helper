package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/mahendraHegde/tradecred-notifier/config"
	conf "github.com/mahendraHegde/tradecred-notifier/config"
	"github.com/mahendraHegde/tradecred-notifier/core"
	"github.com/mahendraHegde/tradecred-notifier/job"
)

func main() {
	config, err := conf.LoadConfig("./")
	if err != nil {
		fmt.Printf("Unable to read config, %v", err)
	}

	r := gin.Default()
	// addCorsRules(r, config)
	wire, err := buildDependencies(config, r)
	if err != nil {
		fmt.Printf("Unable to build dependecies, %v", err)
	}
	r.GET("/live", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, map[string]bool{"live": true})
	})
	core.RegisterCoreRoutes(r, wire.coreController)
	go job.ScheduleDealsCheck(context.Background(), time.Minute*time.Duration(config.TradeCredConfig.DealsCheckSchedule), wire.tradeCredService, wire.callmeBotService)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func addCorsRules(r *gin.Engine, c config.Configurations) {
	allowOrigins := strings.Split(c.Server.Cors, ";")
	fmt.Println("cors>>>", allowOrigins)
	config := cors.DefaultConfig()
	config.AllowOrigins = allowOrigins
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	r.Use(cors.New(config))
}
