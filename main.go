package main

import (
	"fmt"

	"github.com/Kwintenvdb/go-subfinder/client"
	"github.com/spf13/viper"
)

func main() {
	setupConfig()

	apiKey := viper.GetString("apiKey")
	client := client.New(apiKey)
	client.FindSubtitles("Tokyo Vice")
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
