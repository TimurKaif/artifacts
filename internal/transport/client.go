package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/TimurKaif/artifacts/config"
)

type AuthTransport struct {
	Token string
	Base  http.RoundTripper
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.Header.Set("Authorization", "Bearer "+t.Token)
	req2.Header.Set("Accept", "application/json")
	req2.Header.Set("Content-Type", "application/json")
	return t.Base.RoundTrip(req2)
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	character  string
}

type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Transport: &AuthTransport{
				Token: cfg.AuthToken,
				Base:  http.DefaultTransport,
			},
		},
		baseURL:   "https://api.artifactsmmo.com",
		character: cfg.CharacterName,
	}
}

func (c *Client) GetCharacter() (*Response, error) {
	url := fmt.Sprintf("%s/characters/%s", c.baseURL, c.character)
	return c.doRequest("GET", url, nil)
}

func (c *Client) Move(x, y int) (*Response, error) {
	url := fmt.Sprintf("%s/my/%s/action/move", c.baseURL, c.character)

	type MoveRequest struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	jsonData, err := json.Marshal(MoveRequest{X: x, Y: y})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return c.doRequest("POST", url, jsonData)
}

func (c *Client) Attack(participants []string) (*Response, error) {
	url := fmt.Sprintf("%s/my/%s/action/fight", c.baseURL, c.character)

	type AttackRequest struct {
		Participants []string `json:"participants,omitempty"`
	}

	jsonData, err := json.Marshal(AttackRequest{Participants: participants})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return c.doRequest("POST", url, jsonData)
}

func (c *Client) doRequest(method, url string, body []byte) (*Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}
