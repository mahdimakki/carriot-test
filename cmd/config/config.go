package cfg

import (
	"github.com/spf13/viper"
)

func New() (*Config, error) {
	v := viper.New()
	v.SetConfigName("app")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("../../.")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var c Config
	err = v.Unmarshal(&c)
	return &c, err
}

type Config struct {
	MQTT MQTT `mapstructure:"mqtt"`
}

type MQTT struct {
	Host  string `mapstructure:"host"`
	Port  string `mapstructure:"PORT"`
	User  string `mapstructure:"USER"`
	Pass  string `mapstructure:"PASS"`
	Topic string `mapstructure:"TOPIC"`
}
