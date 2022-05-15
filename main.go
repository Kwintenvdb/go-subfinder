package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Kwintenvdb/go-subfinder/client"
	"github.com/Kwintenvdb/go-subfinder/video_finder"
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

		fmt.Println("Trying to find video file in current directory...")
		video, err := videofinder.FindVideoFileInCurrentDir()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found video file: %s\n", video)

		c := createClient()
		c.Login()
		subs := c.FindSubtitles(client.FindSubtitleOptions{
			FileName: video,
			Language: viper.GetString("language"),
		})
		subId := subs.Data[0].Attributes.Files[0].FileId
		c.DownloadSubtitle(subId)
	},
}

func createClient() client.SubtitleClient {
	return client.New(client.ClientConfig{
		ApiKey:   viper.GetString("apiKey"),
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
	})
}

func ExecuteRootCommand() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	cobra.OnInitialize(setupConfig)

	downloadCommand.Flags().StringP("language", "l", "en", "Specify the language")
	viper.BindPFlag("language", downloadCommand.Flags().Lookup("language"))

	rootCommand.AddCommand(downloadCommand)

	ExecuteRootCommand()
}

func setupConfig() {
	executablePath, _ := os.Executable()
	executableConfigPath := filepath.Join(filepath.Dir(executablePath), "/config")

	viper.AddConfigPath("./config")
	viper.AddConfigPath(executableConfigPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
