package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"postfix_to_cf/api"
	"testing"
)

// MockHTTPClient is a mock of HTTPClient
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do is the mock client's implementation of the Do method
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// TestSendPostRequest tests the SendPostRequest function
func TestSendPostRequest(t *testing.T) {
	expectedJSON := `{"to":[{"email":"jane.smith@example.com","name":"Jane Smith"}],"from":{"email":"john.doe@example.com","name":"John Doe"},"subject":"Test Email 1","text":"This is a basic test email."}`

	// Create a mock HTTP client
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Check the request body is as expected
			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)
			if buf.String() != expectedJSON {
				t.Errorf("SendPostRequest body = %v, want %v", buf.String(), expectedJSON)
			}

			// Check the Authorization header is as expected
			authHeader := req.Header.Get("Authorization")
			if authHeader != "exampleToken" {
				t.Errorf("SendPostRequest Authorization header = %v, want %v", authHeader, "exampleToken")
			}

			// Simulate a response
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
			}, nil
		},
	}

	// Load the sample file
	file, err := os.Open("./samples/plain_text.eml")
	if err != nil {
		t.Fatalf("Failed to open sample file: %v", err)
	}
	defer file.Close()

	// Create a request to pass to our handler.
	req := httptest.NewRequest(http.MethodPost, "http://example.com", file)

	// Call the function under test with the mock client
	err = api.SendPostRequest(mockClient, req.URL.String(), "exampleToken", expectedJSON)
	if err != nil {
		t.Errorf("SendPostRequest returned an error: %v", err)
	}
}
