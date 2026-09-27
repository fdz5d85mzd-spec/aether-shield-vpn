package main

import (
 "encoding/json"
 "log"
 "net/http"
 "os"
 "strings"
 "time"
)

type node struct { ID string `json:"id"`; City string `json:"city"`; Country string `json:"country"`; Status string `json:"status"`; LatencyMS *int `json:"latencyMs"`; Capabilities []string `json:"capabilities"`; LastSeen *string `json:"lastSeen"` }
type status struct { Mode string `json:"mode"`; Tunnel string `json:"tunnel"`; API string `json:"api"` }
type sessionRequest struct { Region string `json:"region"`; ClientPublicKey string `json:"clientPublicKey"` }
type sessionResponse struct { Status string `json:"status"`; Reason string `json:"reason"` }

func jsonOut(w http.ResponseWriter, code int, v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(code);_ = json.NewEncoder(w).Encode(v)}

func main(){
 mux:=http.NewServeMux()
 mux.HandleFunc("GET /health",func(w http.ResponseWriter,r *http.Request){jsonOut(w,200,map[string]any{"status":"healthy","time":time.Now().UTC().Format(time.RFC3339)})})
 mux.HandleFunc("GET /v1/system/status",func(w http.ResponseWriter,r *http.Request){jsonOut(w,200,status{Mode:"REAL",Tunnel:"NOT_CONFIGURED",API:"ONLINE"})})
 mux.HandleFunc("GET /v1/nodes",func(w http.ResponseWriter,r *http.Request){jsonOut(w,200,[]node{})})
 mux.HandleFunc("GET /v1/nodes/{id}",func(w http.ResponseWriter,r *http.Request){jsonOut(w,404,map[string]string{"error":"node_not_found","id":r.PathValue("id")})})
 mux.HandleFunc("POST /v1/tunnel/session",func(w http.ResponseWriter,r *http.Request){
  var req sessionRequest
  dec:=json.NewDecoder(http.MaxBytesReader(w,r.Body,16<<10))
  dec.DisallowUnknownFields()
  if err:=dec.Decode(&req);err!=nil{jsonOut(w,400,map[string]string{"error":"invalid_request"});return}
  if strings.TrimSpace(req.ClientPublicKey)==""{jsonOut(w,400,map[string]string{"error":"client_public_key_required"});return}
  // Session provisioning is intentionally unavailable until authenticated node infrastructure exists.
  jsonOut(w,http.StatusServiceUnavailable,sessionResponse{Status:"NOT_CONFIGURED",Reason:"tunnel provisioning infrastructure unavailable"})
 })
 allowed:=strings.TrimSpace(os.Getenv("AETHER_ALLOWED_ORIGINS"))
 handler:=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){origin:=r.Header.Get("Origin");if origin!="" && allowed!="" && origin==allowed{w.Header().Set("Access-Control-Allow-Origin",origin);w.Header().Set("Vary","Origin")} ; mux.ServeHTTP(w,r)})
 srv:=&http.Server{Addr:env("AETHER_API_ADDR",":8080"),Handler:handler,ReadHeaderTimeout:5*time.Second,ReadTimeout:15*time.Second,WriteTimeout:15*time.Second,IdleTimeout:60*time.Second}
 log.Printf("aether control API listening on %s",srv.Addr);log.Fatal(srv.ListenAndServe())
}
func env(k,d string)string{if v:=os.Getenv(k);v!=""{return v};return d}
