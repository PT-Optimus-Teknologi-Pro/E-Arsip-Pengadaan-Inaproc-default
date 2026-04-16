package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"sync-inaproc/config"
	"time"
)

var defaultTransport = &http.Transport{
    Dial: (&net.Dialer{ KeepAlive: 600 * time.Second,}).Dial,
    MaxIdleConns: 20,
    MaxIdleConnsPerHost: 20,
}
var client = http.Client{Transport: defaultTransport, Timeout: 5 * time.Minute}

type Meta struct {
	Cursor		string 		`json:"cursor"`
	HasMore		bool		`json:"has_more"`
	Limit 		int 		`json:"limit"`
}

type ApiResponse struct {
	Success		bool					`json:"success"`
	Data		any						`json:"data"`
	Meta 		Meta					`json:"meta"`
}

func fetchApiCursor(api string, cursor string, model any) error {
	defaultLimit := 100
	url := fmt.Sprintf("%s%s&limit=%d&cursor=%s", DOMAIN_URL, api, defaultLimit, cursor)
	return fetch(url, model)
}


func fetchApi(api string, api_response any) error {
	defaultLimit := 100
	url := fmt.Sprintf("%s%s&limit=%d", DOMAIN_URL, api, defaultLimit)
	return fetch(url, api_response)
}

func fetch(url string, api_response any) error {
	// slog.Info("get request from ", "url", url)
	request, err := http.NewRequest(http.MethodGet, url, nil)
    request.Header.Set("Authorization", config.GetToken())
    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	time.Sleep(time.Duration(config.GetDelay()) * time.Second)
	// slog.Info("delay ", "delays", config.GetDelay(), "request", url)
    response, err := client.Do(request)
	if err != nil {
		slog.Error("Error fetching ", "url", url, "error", err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		slog.Error("HTTP Error", "status", response.Status, "url", url)
		if response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("URL not found (404): %s", url)
		}
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	//
	if err != nil {
		slog.Error("Error reading body ", "url", url, "error", err, "body", buf.String())
		return err
	}
	// slog.Info("response : ", "body", buf.String())
	// var data ApiResponse
	err = json.Unmarshal(buf.Bytes(), api_response)
	if err != nil {
		slog.Error("Failed to unmarshal JSON", "error" , err, "body", buf.String())
		return err
	}
	// slog.Info("response : ", "json", api_response)
	return nil
}

// fetch API use goroutine
func fetchURL(url string, wg *sync.WaitGroup, results chan <- ApiResponse) {
	// Decrement the WaitGroup count when the goroutine finishes
	defer wg.Done()
	slog.Info("get request from "+ url)
	// Make the HTTP request
	request, err := http.NewRequest(http.MethodGet, url, nil)
    request.Header.Set("Authorization", config.GetToken())
    request.Header.Set("Content-Type", "application/json")
    response, err := client.Do(request)
	if err != nil {
		slog.Error("Error fetching %s: %v", url, err)
		results <- ApiResponse{Success: false}
	}
	// Ensure the response body is closed to prevent resource leaks
	defer response.Body.Close()

	// Read the body to ensure connection reuse (optional, but good practice)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("Error reading body for %s: %v", url, err)
		results <- ApiResponse{Success: false}
	}
	// fmt.Println(string(body))
	var data ApiResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		slog.Error("Failed to unmarshal JSON", "error", err)
	}
	results <- data
}
