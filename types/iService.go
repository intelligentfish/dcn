package types

import "github.com/intelligentfish/dcn/define"

// IService service interface
type IService interface {
	// GetServiceType get the service type
	GetServiceType() define.ServiceType
	// GetStartupPriority get the startup priority
	GetStartupPriority() define.StartupPriority
	// GetShutdownPriority get the shutdown priority
	GetShutdownPriority() define.ShutdownPriority
	// Name service name
	Name() string
	// Start start the service
	Start() error
	// Stop stop the service
	Stop()
}
