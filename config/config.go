package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Server          ServerConfig
	Me              MeConfig
	DBConfig        DBConfig
	TradeCredConfig TradeCredConfig
	CallMeBot       CallMeBot
}

//CallMeBot exported
type CallMeBot struct {
	WhatsApp struct {
		ApiKey, Phone, Base string
	}
}

//TradeCredConfig exported
type TradeCredConfig struct {
	Base, RefreshToken string
}

//DBConfig exported
type DBConfig struct {
	Uri, DBName string
}

//MeConfig Exported
type MeConfig struct {
	Email string
}

// ServerConfig exported
type ServerConfig struct {
	Port         int
	Cors, ApiKey string
}

var (
	bindMap = map[string]string{
		"server.port": "PORT",
	}
)

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Configurations, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	//bind env to struct
	for k, v := range bindMap {
		viper.BindEnv(k, v)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// for k, v := range defaults {
	// 	viper.SetDefault(k, v)
	// }

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Config file not found %v", err)
		} else {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
