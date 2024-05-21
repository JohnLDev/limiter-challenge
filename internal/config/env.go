package config

import (
	"encoding/json"
	"os"

	"github.com/spf13/viper"
)

type ConfiguredTokens struct {
	Token string `json:"token"`
	Limit int    `json:"limit"`
}

type config struct {
	RateLimit int    `mapstructure:"RATE_LIMIT"`
	DbPort    int    `mapstructure:"DB_PORT"`
	DbName    string `mapstructure:"DB_NAME"`
	DbHost    string `mapstructure:"DB_HOST"`
	DbPass    string `mapstructure:"DB_PASS"`
	Tokens    []ConfiguredTokens
}

var Conf *config

func loadConfig(path string) (*config, error) {
	var cfg *config
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
func loadTokens() {
	Json, err := os.ReadFile("tokens.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(Json, &Conf.Tokens)
	if err != nil {
		panic(err)
	}
}

func init() {
	cfg, err := loadConfig(".")
	if err != nil {
		panic(err)
	}
	Conf = cfg
	loadTokens()
}
