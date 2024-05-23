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

type Config struct {
	RateLimit int    `mapstructure:"RATE_LIMIT"`
	DbPort    int    `mapstructure:"DB_PORT"`
	DbName    string `mapstructure:"DB_NAME"`
	DbHost    string `mapstructure:"DB_HOST"`
	DbPass    string `mapstructure:"DB_PASS"`
	Tokens    map[string]int
}

func loadConfig(path string) (*Config, error) {
	var cfg *Config
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
func loadTokens(cfg *Config) {
	Json, err := os.ReadFile("tokens.json")
	if err != nil {
		panic(err)
	}
	var tokens []ConfiguredTokens

	err = json.Unmarshal(Json, &tokens)
	if err != nil {
		panic(err)
	}

	cfg.Tokens = make(map[string]int)
	for _, token := range tokens {
		cfg.Tokens[token.Token] = token.Limit
	}
}

func GetConfig() *Config {
	cfg, err := loadConfig(".")
	if err != nil {
		panic(err)
	}
	loadTokens(cfg)
	return cfg
}
