package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiServer = "https://api.opensubtitles.com/api/v1"

type subtitleClient struct {
	clientConfig ClientConfig
}

func New(config ClientConfig) subtitleClient {
	return subtitleClient{
		clientConfig: config,
	}
}

func (c subtitleClient) Login() {
	url := fmt.Sprintf("%s/login", apiServer)

	data := map[string]string{"username": c.clientConfig.Username, "password": c.clientConfig.Password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error marshalling POST body to JSON")
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Api-Key", c.clientConfig.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// fmt.Println(res)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading body: %s", err)
	}

	var loginResponse loginResponseData
	err2 := json.Unmarshal(body, &loginResponse)
	if err2 != nil {
		fmt.Println("error parsing json")
		return
	}
	fmt.Println(loginResponse)
}

type loginResponseData struct {
	User   userData
	Token  string
	Status int
}

type userData struct {
	AllowedDownloads int `json:"allowed_downloads"`
	UserId           int `json:"user_id"`
}

func (c subtitleClient) FindSubtitles(fileName string) {
	url := fmt.Sprintf("%s/subtitles", apiServer)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Api-Key", c.clientConfig.ApiKey)
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
