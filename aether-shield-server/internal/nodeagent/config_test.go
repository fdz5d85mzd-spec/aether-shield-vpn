package nodeagent

import "testing"

func TestLoadConfigFailsClosed(t *testing.T) {
	t.Setenv("AETHER_CONTROL_URL", "")
	t.Setenv("AETHER_NODE_BOOTSTRAP_TOKEN", "")
	t.Setenv("AETHER_NODE_ID", "")
	t.Setenv("AETHER_NODE_ENDPOINT", "")
	t.Setenv("AETHER_NODE_PUBLIC_KEY", "")
	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected missing configuration error")
	}
}

func TestLoadConfigParsesCapabilities(t *testing.T) {
	t.Setenv("AETHER_CONTROL_URL", "https://control.example.invalid")
	t.Setenv("AETHER_NODE_BOOTSTRAP_TOKEN", "runtime-secret")
	t.Setenv("AETHER_NODE_ID", "node-1")
	t.Setenv("AETHER_NODE_ENDPOINT", "vpn.example.invalid:51820")
	t.Setenv("AETHER_NODE_PUBLIC_KEY", "public-only")
	t.Setenv("AETHER_NODE_CAPABILITIES", "wireguard, ipv6 ")
	c, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Capabilities) != 2 {
		t.Fatalf("capabilities=%v", c.Capabilities)
	}
}
