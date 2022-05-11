package main

import (
	"fmt"

	"github.com/Kwintenvdb/go-subfinder/client"
	"github.com/spf13/viper"
)

func main() {
	setupConfig()

	client := client.New(client.ClientConfig{
		ApiKey: viper.GetString("apiKey"),
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
	})
	client.Login()
	// client.FindSubtitles("Tokyo Vice")
}

func setupConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
