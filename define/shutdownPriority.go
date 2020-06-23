package define

// ShutdownPriority shutdown priority
type ShutdownPriority int

const (
	ShutdownPriorityUnknown = ShutdownPriority(iota)
	ShutdownPriorityWeb
	ShutdownPriorityMQ
	ShutdownPriorityDB
)
