package rohr

type HealthInfo struct {
	Verion      string
	Environment string
	Uptime      string
}

type HealthService interface {
	GetHealth() (*HealthInfo, error)
}
