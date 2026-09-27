package nodeagent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Command struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ClientPublicKey string `json:"clientPublicKey"`
	ClientAddress   string `json:"clientAddress"`
}

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewClient(config Config) *Client {
	return &Client{
		baseURL: strings.TrimRight(config.ControlURL, "/"),
		token:   config.Token,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) request(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	return c.http.Do(req)
}

func (c *Client) Heartbeat(ctx context.Context, config Config) error {
	body, err := json.Marshal(map[string]any{
		"id": config.ID, "region": config.Region, "city": config.City, "country": config.Country,
		"endpoint": config.Endpoint, "publicKey": config.PublicKey, "capabilities": config.Capabilities,
	})
	if err != nil {
		return err
	}
	resp, err := c.request(ctx, http.MethodPost, "/v1/nodes/heartbeat", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("heartbeat rejected with status %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) PollCommands(ctx context.Context, nodeID string) ([]Command, error) {
	resp, err := c.request(ctx, http.MethodGet, "/v1/nodes/"+nodeID+"/commands", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("command poll rejected with status %d", resp.StatusCode)
	}
	var commands []Command
	if err := json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&commands); err != nil {
		return nil, err
	}
	return commands, nil
}

func (c *Client) AckCommand(ctx context.Context, nodeID, commandID string, applied bool, message string) error {
	body, err := json.Marshal(map[string]any{"applied": applied, "error": message})
	if err != nil {
		return err
	}
	resp, err := c.request(ctx, http.MethodPost, "/v1/nodes/"+nodeID+"/commands/"+commandID+"/ack", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("command ACK rejected with status %d", resp.StatusCode)
	}
	return nil
}
