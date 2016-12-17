package httprouter_test

import (
	"github.com/concur/rohr"
	rohrHttp "github.com/concur/rohr/http"
	rohrRouter "github.com/concur/rohr/http/httprouter"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockHealthService struct {
	*rohr.HealthInfo
}

func (h *mockHealthService) GetHealth() (*rohr.HealthInfo, error) {
	return h.HealthInfo, nil
}

func TestRouter_GetHealthHandler(t *testing.T) {
	healthHandler := rohrRouter.GetHealthHandler(&mockHealthService{
		HealthInfo: &rohr.HealthInfo{
			Verion:      "testing",
			Environment: "dev",
			Uptime:      "1s",
		},
	})

	router := httprouter.New()
	router.GET("/health_test", healthHandler)
	r, _ := http.NewRequest(http.MethodGet, "/health_test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	response_body := w.Body.String()
	response_header := w.Header()

	if w.Code != http.StatusOK {
		t.Errorf("GetHealthHandler returns HTTP error code: %v, header: %#v, body: %#v", w.Code, response_header, response_body)
	}

	if !strings.Contains(response_body, "{\"Verion\":\"testing\",\"Environment\":\"dev\",\"Uptime\":\"1s\"}") {
		t.Errorf("GetHealthHandler should return mock health info.")
	}

	t.Logf("GetHealthHandler returns header: %#v, body: %#v", response_header, response_body)
}

func TestRouter_RegisterRoute(t *testing.T) {
	router := rohrRouter.NewRouter()
	hrouter := router.RegisterRoute(&rohrHttp.ApiServer{})
	r, _ := http.NewRequest(http.MethodGet, rohrRouter.HEALTH_PATH, nil)
	w := httptest.NewRecorder()
	hrouter.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("RegisterRoute should add Health endpoint hanlder for \"/health\" path. Request returns HTTP error code: %v, header: %#v, body: %#v", w.Code, w.Header(), w.Body.String())
	}
}
