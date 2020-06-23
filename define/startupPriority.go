package define

// StartupPriority startup priority
type StartupPriority int

const (
	StartupPriorityUnknown = StartupPriority(iota)
	StartupPriorityDB
	StartupPriorityMQ
	StartupPriorityWeb
)
