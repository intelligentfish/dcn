package define

// ServiceType service type
type ServiceType int

const (
	ServiceTypeUnknown = ServiceType(iota)
	ServiceTypeDB
	ServiceTypeMQ
	ServiceTypeWeb
)
