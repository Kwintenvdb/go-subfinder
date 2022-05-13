package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const apiServer = "https://api.opensubtitles.com/api/v1"

type subtitleClient struct {
	clientConfig ClientConfig
	loginData    *loginResponseData
}

func New(config ClientConfig) subtitleClient {
	return subtitleClient{
		clientConfig: config,
	}
}

func (c subtitleClient) Login() {
	fmt.Println("Logging in...")

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

func (c subtitleClient) FindSubtitles(fileName string) QueryResponse {
	fmt.Printf("Finding subtitles for file %s...\n", fileName)

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

	return fromJson[QueryResponse](body)
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

func (c subtitleClient) DownloadSubtitle(fileId int) {
	data := map[string]int{"file_id": fileId}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error marshalling POST body to JSON")
		return
	}

	url := fmt.Sprintf("%s/download", apiServer)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Api-Key", c.clientConfig.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading body: %s", err)
	}
	downloadData := fromJson[DownloadResponse](body)
	fmt.Printf("Successfully retrieved download link. Remaining downloads: %d\n", downloadData.Remaining)

	out, err := os.Create(downloadData.FileName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	downloadRes, err := client.Get(downloadData.Link)
	if err != nil {
		panic(err)
	}
	defer downloadRes.Body.Close()

	io.Copy(out, downloadRes.Body)
	fmt.Printf("Successfully downloaded subtitle file: %s\n", downloadData.FileName)
}

type DownloadResponse struct {
	Link      string
	FileName  string `json:"file_name"`
	Remaining int
	Message   string
}
