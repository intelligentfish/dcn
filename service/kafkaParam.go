package service

import "encoding/json"

// KafkaParam kafka parameter
type KafkaParam struct {
	Network   string `json:"network"`   // kafka network
	Address   string `json:"address"`   // kafka address
	Topic     string `json:"topic"`     // kafka topic
	Partition int    `json:"partition"` // kafka partition
	Offset    int64  `json:"offset"`    // kafka offset
	Whence    *int   `json:"whence"`    // kafka whence
}

// String string description
func (object *KafkaParam) String() string {
	raw, _ := json.Marshal(object)
	return string(raw)
}
