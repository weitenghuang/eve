package service_test

import (
	"github.com/concur/rohr"
	"github.com/concur/rohr/pkg/config"
	"github.com/concur/rohr/service"
	"reflect"
	"testing"
)

var healthService *service.HealthService

func TestHealthService_New(t *testing.T) {
	healthService = service.NewHealthService()
	if !reflect.DeepEqual(healthService, &service.HealthService{}) {
		t.Errorf("Unexpected HealthService value: %#v", healthService)
	}
}

func TestHealthService_GetHealth(t *testing.T) {
	errCount := 6
	sysConfig := config.NewSystemConfig()
	mock_healthInfo := &rohr.HealthInfo{
		Hostname: sysConfig.Hostname,
		Metadata: map[string]string{
			"Version":     sysConfig.Version,
			"Environment": sysConfig.Environment,
		},
		Errors: make([]rohr.Error, errCount),
		Uptime: "",
	}

	healthInfo := healthService.GetHealth()

	if len(healthInfo.Uptime) == len(mock_healthInfo.Uptime) {
		t.Errorf("HealthInfo Uptime should not be empty string. Uptime: %v", healthInfo.Uptime)
	}
	if len(healthInfo.Errors) != len(mock_healthInfo.Errors) {
		t.Errorf("HealthInfo should report %v errors when none of dependencies is available. Errors: %#v", errCount, healthInfo.Errors)
	}
	if !reflect.DeepEqual(healthInfo.Metadata, mock_healthInfo.Metadata) {
		t.Errorf("Unexpected healthInfo metadata value: %#v", healthInfo.Metadata)
	}
}
