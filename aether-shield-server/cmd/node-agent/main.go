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
	peers := nodeagent.PeerManager{Interface: config.WireGuardInterface, Timeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cycle := func() {
		probeCtx, probeCancel := context.WithTimeout(ctx, 6*time.Second)
		defer probeCancel()
		if err := runtime.Verify(probeCtx, config); err != nil {
			log.Printf("runtime verification failed; heartbeat/commands suppressed: %v", err)
			return
		}
		hbCtx, hbCancel := context.WithTimeout(ctx, 10*time.Second)
		if err := client.Heartbeat(hbCtx, config); err != nil {
			hbCancel()
			log.Printf("heartbeat failed: %v", err)
			return
		}
		hbCancel()

		pollCtx, pollCancel := context.WithTimeout(ctx, 10*time.Second)
		commands, err := client.PollCommands(pollCtx, config.ID)
		pollCancel()
		if err != nil {
			log.Printf("command poll failed: %v", err)
			return
		}
		for _, command := range commands {
			cmdCtx, cmdCancel := context.WithTimeout(ctx, 6*time.Second)
			var applyErr error
			switch command.Type {
			case "ADD_PEER":
				applyErr = peers.Add(cmdCtx, command.ClientPublicKey, command.ClientAddress)
			case "REMOVE_PEER":
				applyErr = peers.Remove(cmdCtx, command.ClientPublicKey)
			default:
				applyErr = &unsupportedCommandError{command.Type}
			}
			cmdCancel()

			message := ""
			if applyErr != nil {
				message = applyErr.Error()
			}
			ackCtx, ackCancel := context.WithTimeout(ctx, 10*time.Second)
			if err := client.AckCommand(ackCtx, config.ID, command.ID, applyErr == nil, message); err != nil {
				log.Printf("command %s ACK failed: %v", command.ID, err)
			}
			ackCancel()
		}
	}

	cycle()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cycle()
		}
	}
}

type unsupportedCommandError struct{ commandType string }

func (e *unsupportedCommandError) Error() string { return "unsupported command type: " + e.commandType }
