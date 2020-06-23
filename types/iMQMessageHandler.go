package types

// IMQMessageHandler mq message handler
type IMQMessageHandler interface {
	// OnMQMessage message processor
	OnMQMessage(raw []byte, offset int64)
}
