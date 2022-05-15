package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	apiBaseUrl   = "https://api.opensubtitles.com/api/v1"
	loginUrl     = apiBaseUrl + "/login"
	subtitlesUrl = apiBaseUrl + "/subtitles"
	downloadUrl  = apiBaseUrl + "/download"
)

type SubtitleClient struct {
	clientConfig ClientConfig
	loginData    *loginResponseData
}

func New(config ClientConfig) SubtitleClient {
	return SubtitleClient{
		clientConfig: config,
	}
}

func (c SubtitleClient) Login() {
	data := map[string]string{"username": c.clientConfig.Username, "password": c.clientConfig.Password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error marshalling POST body to JSON")
		return
	}

	fmt.Println("Logging in...")
	req, err := http.NewRequest("POST", loginUrl, bytes.NewBuffer(jsonData))
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading body: %s", err)
	}

	loginResponse := fromJson[loginResponseData](body)
	c.loginData = &loginResponse

	fmt.Println("Logged in successfully.")
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

type FindSubtitleOptions struct {
	FileName string
	Language string
}

func (c SubtitleClient) FindSubtitles(options FindSubtitleOptions) (QueryResponse, error) {
	fmt.Printf("Finding subtitles for file %s in language %s...\n", options.FileName, options.Language)

	req, err := http.NewRequest("GET", subtitlesUrl, nil)
	if err != nil {
		return QueryResponse{}, err
	}
	req.Header.Add("Api-Key", c.clientConfig.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	query := req.URL.Query()
	query.Add("query", options.FileName)
	query.Add("languages", options.Language)
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("error: %s", err)
		return QueryResponse{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading body: %s", err)
		return QueryResponse{}, err
	}

	if res.StatusCode != 200 {
		fmt.Printf("Received status code %d, response body: %s\n", res.StatusCode, string(body))
		return QueryResponse{}, errors.New("status code of subtitle search response was not 200")
	}

	return fromJson[QueryResponse](body), nil
}

func fromJson[T any](responseBody []byte) T {
	var data T
	err := json.Unmarshal(responseBody, &data)
	if err != nil {
		panic("error parsing json")
	}
	return data
}

type QueryResponse struct {
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
	Page       int
	Data       []QueryResponseData
}

type QueryResponseData struct {
	Id         string
	DataType   string `json:"type"`
	Attributes QueryResponseDataAttributes
}

type QueryResponseDataAttributes struct {
	SubtitleId    string `json:"subtitle_id"`
	Language      string
	DownloadCount int `json:"download_count"`
	Url           string
	Files         []FileData
}

type FileData struct {
	FileId   int    `json:"file_id"`
	FileName string `json:"file_name"`
}

func (c SubtitleClient) DownloadSubtitle(fileId int) error {
	data := map[string]int{"file_id": fileId}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error marshalling POST body to JSON")
		return err
	}

	req, err := http.NewRequest("POST", downloadUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Api-Key", c.clientConfig.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading body: %s", err)
		return err
	}
	downloadData := fromJson[DownloadResponse](body)
	fmt.Printf("Successfully retrieved download link. Remaining downloads: %d\n", downloadData.Remaining)

	out, err := os.Create(downloadData.FileName)
	if err != nil {
		return err
	}
	defer out.Close()

	downloadRes, err := client.Get(downloadData.Link)
	if err != nil {
		return err
	}
	defer downloadRes.Body.Close()

	io.Copy(out, downloadRes.Body)
	fmt.Printf("Successfully downloaded subtitle file: %s\n", downloadData.FileName)
	return nil
}

type DownloadResponse struct {
	Link      string
	FileName  string `json:"file_name"`
	Remaining int
	Message   string
}
