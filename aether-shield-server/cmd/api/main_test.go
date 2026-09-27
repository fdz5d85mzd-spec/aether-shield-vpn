package main
import("net/http/httptest";"testing")
func TestJSONOut(t *testing.T){w:=httptest.NewRecorder();jsonOut(w,200,map[string]string{"status":"healthy"});if w.Code!=200{t.Fatalf("status=%d",w.Code)};if got:=w.Header().Get("Content-Type");got!="application/json"{t.Fatalf("content-type=%s",got)}}
func TestEnvDefault(t *testing.T){if got:=env("AETHER_TEST_MISSING","fallback");got!="fallback"{t.Fatalf("got %q",got)}}
