package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// DeviceCodeResponse from the pairing service.
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// DeviceTokenResponse from the pairing service on successful exchange.
type DeviceTokenResponse struct {
	AccessToken string `json:"access_token"`
	APIKey      string `json:"api_key"`
}

// DeviceClient communicates with the pairing service.
type DeviceClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewDeviceClient creates a client for the pairing service.
func NewDeviceClient(baseURL string) *DeviceClient {
	return &DeviceClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// RequestCode initiates a device authorization flow.
func (c *DeviceClient) RequestCode() (*DeviceCodeResponse, error) {
	body := bytes.NewBufferString(`{"client_id":"trello-cli"}`)
	resp, err := c.httpClient.Post(c.baseURL+"/device/code", "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("requesting device code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pairing service returned %d", resp.StatusCode)
	}

	var result DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding device code response: %w", err)
	}

	return &result, nil
}

// PollToken makes a single poll request for the token.
func (c *DeviceClient) PollToken(deviceCode string) (*DeviceTokenResponse, error) {
	payload := map[string]string{
		"device_code": deviceCode,
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
	}
	body, _ := json.Marshal(payload)

	resp, err := c.httpClient.Post(c.baseURL+"/token", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("polling for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result DeviceTokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("decoding token response: %w", err)
		}
		return &result, nil
	}

	var errResp struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return nil, fmt.Errorf("pairing service returned %d", resp.StatusCode)
	}

	return nil, errors.New(errResp.Error)
}
