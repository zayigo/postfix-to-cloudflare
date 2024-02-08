package api

import (
	"bytes"
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// SendPostRequest sends a POST request to the specified endpoint with the given token and body data
func SendPostRequest(client HTTPClient, endpoint, token, bodyData string) error {
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(bodyData)))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status code: %d", resp.StatusCode)
	}

	return nil
}
