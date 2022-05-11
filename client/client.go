package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiServer = "https://api.opensubtitles.com/api/v1"

type subtitleClient struct {
	apiKey string
}

func New(apiKey string) subtitleClient {
	return subtitleClient{
		apiKey: apiKey,
	}
}

func (c subtitleClient) FindSubtitles(fileName string) {
	url := fmt.Sprintf("%s/subtitles", apiServer)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Api-Key", c.apiKey)
	req.Header.Add("Content-Type", "application/json")

	query := req.URL.Query()
	query.Add("query", fileName)
	req.URL.RawQuery = query.Encode()


	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("error: %s", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading body: %s", err)
	}

	toJson(body)
}

func toJson(responseBody []byte) {
	var data queryResponse
	err := json.Unmarshal(responseBody, &data)
	if err != nil {
		fmt.Println("error parsing json")
		return
	}
	fmt.Println(data)
}

type queryResponse struct {
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
	Page       int
	Data       []queryResponseData
}

type queryResponseData struct {
	Id         string
	DataType   string `json:"type"`
	Attributes queryResponseDataAttributes
}

type queryResponseDataAttributes struct {
	SubtitleId    string `json:"subtitle_id"`
	Language      string
	DownloadCount int `json:"download_count"`
	Url           string
}
