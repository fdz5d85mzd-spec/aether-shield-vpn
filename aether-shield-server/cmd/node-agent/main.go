package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/fdz5d85mzd-spec/aether-shield-vpn/aether-shield-server/internal/nodeagent"
)

func main() {
	config, err := nodeagent.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	client := nodeagent.NewClient(config)
	runtime := nodeagent.WireGuardRuntime{Interface: config.WireGuardInterface, Timeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	send := func() {
		probeCtx, probeCancel := context.WithTimeout(ctx, 6*time.Second)
		defer probeCancel()
		if err := runtime.Verify(probeCtx, config); err != nil {
			log.Printf("runtime verification failed; heartbeat suppressed: %v", err)
			return
		}
		hbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := client.Heartbeat(hbCtx, config); err != nil {
			log.Printf("heartbeat failed: %v", err)
		}
	}

	send()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			send()
		}
	}
}
