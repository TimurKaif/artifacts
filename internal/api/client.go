package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"artifacts/config"
)

type Client struct {
	config     config.Config
	httpClient *http.Client
	baseURL    string
}

type CharacterInfo struct {
	Data struct {
		Name string `json:"name"`
		X    int    `json:"x"`
		Y    int    `json:"y"`
		HP   int    `json:"hp"`
	} `json:"data"`
}

type AttackRequest struct {
	Participants []string `json:"participants,omitempty"`
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		config:     cfg,
		httpClient: &http.Client{},
		baseURL:    "https://api.artifactsmmo.com",
	}
}

func (c *Client) doRequest(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.AuthToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) GetCharacter() (*CharacterInfo, error) {
	url := fmt.Sprintf("%s/characters/%s", c.baseURL, c.config.CharacterName)

	body, err := c.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var info CharacterInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &info, nil
}

func (c *Client) Move(x, y int) (string, error) {
	url := fmt.Sprintf("%s/my/%s/action/move", c.baseURL, c.config.CharacterName)

	type MoveRequest struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	jsonData, err := json.Marshal(MoveRequest{X: x, Y: y})
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	body, err := c.doRequest("POST", url, jsonData)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) Attack(participants []string) (string, error) {
	url := fmt.Sprintf("%s/my/%s/action/fight", c.baseURL, c.config.CharacterName)

	payload := AttackRequest{
		Participants: participants,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	body, err := c.doRequest("POST", url, jsonData)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
