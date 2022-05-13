package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Kwintenvdb/go-subfinder/client"
	mediafinder "github.com/Kwintenvdb/go-subfinder/media_finder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCommand = &cobra.Command{
	Use: "go-subfinder",
	Run: func(cmd *cobra.Command, args []string) {
		println("executing command")
	},
}

var downloadCommand = &cobra.Command{
	Use: "download",
	Run: func(cmd *cobra.Command, args []string) {
		println("executing download command")

		println("Listing files in dir")
		video, err := mediafinder.FindVideoFileInCurrentDir()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found video file: %s\n", video)

		client := client.New(client.ClientConfig{
			ApiKey:   viper.GetString("apiKey"),
			Username: viper.GetString("username"),
			Password: viper.GetString("password"),
		})
		client.Login()
		subs := client.FindSubtitles(video)
		subId := subs.Data[0].Attributes.Files[0].FileId
		client.DownloadSubtitle(subId)
	},
}

func ExecuteRootCommand() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	rootCommand.AddCommand(downloadCommand)

	cobra.OnInitialize(setupConfig)
	ExecuteRootCommand()
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
