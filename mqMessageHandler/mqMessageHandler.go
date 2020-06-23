package mqMessageHandler

import "github.com/intelligentfish/dcn/types"

// MQMessageHandler mq message handler
type MQMessageHandler struct {
	callback types.MQMessageCallback
}

// NewMQMessageHandler factory method
func NewMQMessageHandler(callback types.MQMessageCallback) *MQMessageHandler {
	return &MQMessageHandler{callback: callback}
}

// OnMQMessage process mq message
func (object *MQMessageHandler) OnMQMessage(raw []byte, offset int64) {
}
