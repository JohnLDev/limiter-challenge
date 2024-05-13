package config

import "github.com/spf13/viper"

type config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Port        int    `mapstructure:"PORT"`
}

var Conf *config

func loadConfig(path string) (*config, error) {
	var cfg *config
	viper.SetDefault("PORT", 3001)
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

func init() {
	cfg, err := loadConfig(".")
	if err != nil {
		panic(err)
	}
	Conf = cfg
}
