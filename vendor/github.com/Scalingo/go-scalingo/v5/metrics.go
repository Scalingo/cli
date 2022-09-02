package scalingo

// These constants contain the string repsentation for all metric available on Scalingo
const (
	MetricMemory                = "memory"
	MetricSwap                  = "swap"
	MetricCPU                   = "cpu"
	MetricRouter5XX             = "5XX"
	MetricRouterAll             = "all"
	MetricRouterServersAmount   = "servers_amount"
	MetricRouterRPMPerContainer = "rpm_per_container"
	MetricRouterP95ResponseTime = "p95_response_time"
)
