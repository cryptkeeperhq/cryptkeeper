package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
)

type APIActionPlugin struct {
	// You can add any configuration or dependencies needed for the plugin here
}

func (a *APIActionPlugin) Execute(inputs []models.NodeInput) (map[string]interface{}, error) {
	// Define the API endpoint URL
	// apiURL := "https://httpbin.org" // Replace with your API URL

	// Initialize variables to store HTTP method and request body
	var apiURL string
	var method string
	var requestBodyJSON []byte

	// Extract inputs and populate the HTTP method and request body
	for _, input := range inputs {
		if input.Value == nil {
			continue
		}
		switch input.ID {
		case "url":
			apiURL = input.Value.(string)
		case "method":
			method = input.Value.(string)
		case "data":
			if input.Value.(string) == "" {
				continue
			}
			if err := json.Unmarshal([]byte(input.Value.(string)), &input.Value); err != nil {
				return nil, fmt.Errorf("error unmarshaling input.Value: %v", err)
			}

			data := input.Value.(map[string]interface{})
			requestBodyJSON, _ = json.Marshal(data)
		}
	}

	// Validate that a valid HTTP method is provided
	if method != "GET" && method != "POST" {
		return nil, fmt.Errorf("invalid HTTP method: %s", method)
	}

	// Create an HTTP client
	client := createHTTPClient()

	// Send the HTTP request based on the specified method
	resp, err := sendHTTPRequest(client, method, apiURL, requestBodyJSON)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Parse the response JSON
	result, err := parseAPIResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	// Return the API response or error
	return result, nil
}

func createHTTPClient() *http.Client {
	return &http.Client{}
}

func sendHTTPRequest(client *http.Client, method, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set request headers if needed
	// req.Header.Set("Content-Type", "application/json")
	// Add any other headers as necessary

	// Send the HTTP request
	return client.Do(req)
}

func parseAPIResponse(body io.Reader) (map[string]interface{}, error) {
	// Read the response body
	responseBody, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	// Parse the response JSON
	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
