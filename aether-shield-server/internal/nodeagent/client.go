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

func (c *Client) Heartbeat(ctx context.Context, config Config) error {
	body, err := json.Marshal(map[string]any{
		"id":           config.ID,
		"region":       config.Region,
		"city":         config.City,
		"country":      config.Country,
		"endpoint":     config.Endpoint,
		"publicKey":    config.PublicKey,
		"capabilities": config.Capabilities,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/nodes/heartbeat", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
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
