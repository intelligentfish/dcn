package service

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"testing"
	"time"
)

func TestMQSrv(t *testing.T) {
	defer func() { MQSrvInst().Stop() }()
	err := MQSrvInst(KafkaParamOption(&KafkaParam{
		Network:   "tcp",
		Address:   "localhost:9092",
		Topic:     "test1",
		Partition: 0,
	})).Start()
	if nil != err {
		t.Error(err)
		return
	}
	for i := 0; i < 1024; i++ {
		_, err = MQSrvInst().Write("test1", []byte(fmt.Sprintf("value is: %d", i)))
		if nil != err {
			t.Error(err)
			return
		}
	}
	var msg kafka.Message
	offset := int64(0)
	whence := -1
	for i := 0; i < 102400; i++ {
		msg, err = MQSrvInst().Read("test1", offset, whence, 4096, time.Now().Add(time.Second))
		if nil != err {
			if e, ok := err.(kafka.Error); ok && e.Timeout() {
				break
			}
			t.Error(err)
			return
		}
		fmt.Println(msg, string(msg.Value))
	}
}
