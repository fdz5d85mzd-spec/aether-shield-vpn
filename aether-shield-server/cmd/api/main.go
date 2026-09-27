package main

import (
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fdz5d85mzd-spec/aether-shield-vpn/aether-shield-server/internal/nodes"
	"github.com/fdz5d85mzd-spec/aether-shield-vpn/aether-shield-server/internal/sessions"
)

type systemStatus struct {
	Mode   string `json:"mode"`
	Tunnel string `json:"tunnel"`
	API    string `json:"api"`
}

type sessionRequest struct {
	Region          string `json:"region"`
	ClientPublicKey string `json:"clientPublicKey"`
}

type heartbeatRequest struct {
	ID           string   `json:"id"`
	Region       string   `json:"region"`
	City         string   `json:"city"`
	Country      string   `json:"country"`
	Endpoint     string   `json:"endpoint"`
	PublicKey    string   `json:"publicKey"`
	Capabilities []string `json:"capabilities"`
}

func jsonOut(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func authorizedNode(r *http.Request) bool {
	expected := os.Getenv("AETHER_NODE_BOOTSTRAP_TOKEN")
	provided := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if expected == "" || provided == "" || len(expected) != len(provided) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expected), []byte(provided)) == 1
}

func main() {
	registry := nodes.NewRegistry()
	provisioner := sessions.NewProvisioner(registry)
	commands := nodes.NewCommandQueue()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		jsonOut(w, http.StatusOK, map[string]any{"status": "healthy", "time": time.Now().UTC().Format(time.RFC3339)})
	})
	mux.HandleFunc("GET /v1/system/status", func(w http.ResponseWriter, r *http.Request) {
		jsonOut(w, http.StatusOK, systemStatus{Mode: "REAL", Tunnel: "NOT_CONFIGURED", API: "ONLINE"})
	})
	mux.HandleFunc("GET /v1/nodes", func(w http.ResponseWriter, r *http.Request) {
		registry.MarkStaleOffline(time.Now().UTC(), 90*time.Second)
		jsonOut(w, http.StatusOK, registry.List())
	})
	mux.HandleFunc("GET /v1/nodes/{id}", func(w http.ResponseWriter, r *http.Request) {
		registry.MarkStaleOffline(time.Now().UTC(), 90*time.Second)
		n, ok := registry.Get(r.PathValue("id"))
		if !ok {
			jsonOut(w, http.StatusNotFound, map[string]string{"error": "node_not_found"})
			return
		}
		jsonOut(w, http.StatusOK, n)
	})
	mux.HandleFunc("POST /v1/nodes/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		if !authorizedNode(r) {
			jsonOut(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}
		var req heartbeatRequest
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
			jsonOut(w, http.StatusBadRequest, map[string]string{"error": "invalid_request"})
			return
		}
		n := nodes.Node{ID: req.ID, Region: req.Region, City: req.City, Country: req.Country, Endpoint: req.Endpoint, PublicKey: req.PublicKey, Status: nodes.StatusUnknown, Capabilities: req.Capabilities}
		if err := registry.Upsert(n); err != nil {
			jsonOut(w, http.StatusBadRequest, map[string]string{"error": "invalid_node"})
			return
		}
		now := time.Now().UTC()
		_ = registry.Heartbeat(req.ID, now)
		jsonOut(w, http.StatusOK, map[string]any{"status": "ONLINE", "lastSeen": now})
	})
	mux.HandleFunc("GET /v1/nodes/{id}/commands", func(w http.ResponseWriter, r *http.Request) {
		if !authorizedNode(r) {
			jsonOut(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}
		nodeID := r.PathValue("id")
		if _, ok := registry.Get(nodeID); !ok {
			jsonOut(w, http.StatusNotFound, map[string]string{"error": "node_not_found"})
			return
		}
		jsonOut(w, http.StatusOK, commands.Poll(nodeID))
	})
	mux.HandleFunc("POST /v1/nodes/{id}/commands/{commandId}/ack", func(w http.ResponseWriter, r *http.Request) {
		if !authorizedNode(r) {
			jsonOut(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}
		var ack struct {
			Applied bool   `json:"applied"`
			Error   string `json:"error"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&ack); err != nil {
			jsonOut(w, http.StatusBadRequest, map[string]string{"error": "invalid_request"})
			return
		}
		cmd, err := commands.Ack(r.PathValue("id"), r.PathValue("commandId"), ack.Applied, ack.Error)
		if err != nil {
			jsonOut(w, http.StatusNotFound, map[string]string{"error": "command_not_found"})
			return
		}
		jsonOut(w, http.StatusOK, cmd)
	})
	mux.HandleFunc("POST /v1/tunnel/session", func(w http.ResponseWriter, r *http.Request) {
		var req sessionRequest
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
			jsonOut(w, http.StatusBadRequest, map[string]string{"error": "invalid_request"})
			return
		}
		registry.MarkStaleOffline(time.Now().UTC(), 90*time.Second)
		p, err := provisioner.Provision(sessions.Request{Region: req.Region, ClientPublicKey: req.ClientPublicKey})
		if err != nil {
			jsonOut(w, http.StatusServiceUnavailable, map[string]string{"status": "NOT_CONFIGURED", "reason": err.Error()})
			return
		}
		cmd, err := commands.Enqueue(p.NodeID, nodes.Command{
			Type: nodes.CommandAddPeer, ClientPublicKey: req.ClientPublicKey, ClientAddress: p.ClientAddress,
		})
		if err != nil {
			jsonOut(w, http.StatusServiceUnavailable, map[string]string{"status": "COMMAND_QUEUE_FAILED"})
			return
		}
		jsonOut(w, http.StatusAccepted, map[string]any{
			"nodeId": p.NodeID, "endpoint": p.Endpoint, "serverPublicKey": p.ServerPublicKey,
			"clientAddress": p.ClientAddress, "status": "PENDING_NODE_APPLY", "commandId": cmd.ID,
		})
	})

	allowed := strings.TrimSpace(os.Getenv("AETHER_ALLOWED_ORIGINS"))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && allowed != "" && origin == allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		mux.ServeHTTP(w, r)
	})
	srv := &http.Server{Addr: env("AETHER_API_ADDR", ":8080"), Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	log.Printf("aether control API listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
