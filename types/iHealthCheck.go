package types

// IHealthCheck health check
type IHealthCheck interface {
	OnHealthCheck() (err error)
}
