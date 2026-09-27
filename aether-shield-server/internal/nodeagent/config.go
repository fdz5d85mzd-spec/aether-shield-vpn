package nodeagent

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	ControlURL         string
	Token              string
	ID                 string
	Region             string
	City               string
	Country            string
	Endpoint           string
	PublicKey          string
	WireGuardInterface string
	Capabilities       []string
}

func LoadConfig() (Config, error) {
	c := Config{
		ControlURL:         os.Getenv("AETHER_CONTROL_URL"),
		Token:              os.Getenv("AETHER_NODE_BOOTSTRAP_TOKEN"),
		ID:                 os.Getenv("AETHER_NODE_ID"),
		Region:             os.Getenv("AETHER_NODE_REGION"),
		City:               os.Getenv("AETHER_NODE_CITY"),
		Country:            os.Getenv("AETHER_NODE_COUNTRY"),
		Endpoint:           os.Getenv("AETHER_NODE_ENDPOINT"),
		PublicKey:          os.Getenv("AETHER_NODE_PUBLIC_KEY"),
		WireGuardInterface: os.Getenv("AETHER_WG_INTERFACE"),
	}
	if raw := strings.TrimSpace(os.Getenv("AETHER_NODE_CAPABILITIES")); raw != "" {
		for _, value := range strings.Split(raw, ",") {
			if value = strings.TrimSpace(value); value != "" {
				c.Capabilities = append(c.Capabilities, value)
			}
		}
	}
	if c.ControlURL == "" || c.Token == "" || c.ID == "" || c.Endpoint == "" || c.PublicKey == "" || c.WireGuardInterface == "" {
		return Config{}, errors.New("control URL, token, node ID, endpoint, public key and WireGuard interface are required")
	}
	return c, nil
}
