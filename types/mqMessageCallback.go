package types

// MQMessageCallback mq message callback
type MQMessageCallback func(raw []byte, offset int)
