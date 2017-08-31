package httprouter

import (
	"github.com/julienschmidt/httprouter"
	"github.com/scipian/eve"
	eveHttp "github.com/scipian/eve/http"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockHealthService struct {
}

func (h *mockHealthService) GetHealth() *eve.HealthInfo {
	return &eve.HealthInfo{
		Hostname: "TestHost",
		Metadata: map[string]string{
			"Version":     "v0.0.1",
			"Environment": "DEV",
		},
		Uptime: "",
	}
}

func TestRouter_getHealthHandler(t *testing.T) {
	healthHandler := getHealthHandler(&mockHealthService{})

	router := httprouter.New()
	router.GET("/health_test", healthHandler)
	r, _ := http.NewRequest(http.MethodGet, "/health_test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	response_body := w.Body.String()
	response_header := w.Header()

	if w.Code != http.StatusOK {
		t.Errorf("getHealthHandler returns HTTP error code: %v, header: %#v, body: %#v", w.Code, response_header, response_body)
	}

	if !strings.Contains(response_body, "\"metadata\":{\"Environment\":\"DEV\",\"Version\":\"v0.0.1\"}") {
		t.Errorf("getHealthHandler should return mock health info.")
	}
}

func TestRouter_RegisterRoute(t *testing.T) {
	router := NewRouter()
	hrouter := router.RegisterRoute(&eveHttp.ApiServer{})
	r, _ := http.NewRequest(http.MethodGet, HEALTH_PATH, nil)
	w := httptest.NewRecorder()
	hrouter.ServeHTTP(w, r)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("RegisterRoute should add Health endpoint hanlder for \"/health\" path. Without dependencies, request should returns HTTP error code 503. Return code: %v, header: %#v, body: %#v", w.Code, w.Header(), w.Body.String())
	}
}
