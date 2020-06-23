package srv

import (
	"context"
	"github.com/intelligentfish/dcn/define"
	"github.com/intelligentfish/dcn/log"
	"github.com/intelligentfish/dcn/mqMessageHandler"
	"github.com/intelligentfish/dcn/srvGroup"
	"github.com/intelligentfish/dcn/types"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	defaultReadTimeout  = 5 * time.Second
	defaultReadMaxBytes = 4096
)

var (
	mqSrvInst *MQSrv    // kafka service instance
	mqSrvOnce sync.Once // kafka service once
)

// MQSrv mq service
type MQSrv struct {
	*BaseSrv
	kafkaParamList      []*KafkaParam                        // kafka param list
	kafkaConnMapMutex   sync.RWMutex                         // kafka topic to connection map read write lock
	kafkaConnMap        map[string]*kafka.Conn               // kafka topic to connection map
	subscribeGroupMutex sync.RWMutex                         // kafka message subscribe group read write lock
	subscribeGroup      map[string][]types.IMQMessageHandler // kafka message subscribe group
}

// NewMQSrv factory method
func NewMQSrv(options ...MQSrvOption) *MQSrv {
	object := &MQSrv{
		BaseSrv:        NewSrvBase("MQSrv", define.ServiceTypeMQ, define.StartupPriorityMQ, define.ShutdownPriorityMQ),
		kafkaConnMap:   make(map[string]*kafka.Conn, 0),
		subscribeGroup: make(map[string][]types.IMQMessageHandler, 0),
	}
	object.BaseSrv.Runner = object
	for _, opt := range options {
		opt(object)
	}
	return object
}

// dial dial kafka connection
func (object *MQSrv) dial() (err error) {
	object.kafkaConnMapMutex.Lock()
	defer object.kafkaConnMapMutex.Unlock()
	var conn *kafka.Conn
	for _, param := range object.kafkaParamList {
		log.Inst().Info("dial kafka connection", zap.String("param", param.String()))
		if conn, err = kafka.DialLeader(context.Background(),
			param.Network,
			param.Address,
			param.Topic,
			param.Partition); nil != err {
			return
		}
		if nil != param.Whence {
			if _, err = conn.Seek(param.Offset, *param.Whence); nil != err {
				return
			}
		}
		object.kafkaConnMap[param.Topic] = conn
	}
	return
}

// clean clean resource
func (object *MQSrv) clean() {
	object.kafkaConnMapMutex.Lock()
	defer object.kafkaConnMapMutex.Unlock()
	var err error
	for topic, conn := range object.kafkaConnMap {
		log.Inst().Info("close kafka connection", zap.String("topic", topic))
		if err = conn.Close(); nil != err {
			log.Inst().Error("close kafka connection", zap.String("error", err.Error()))
		}
	}
}

// Publish publish mq message
func (object *MQSrv) Publish(topic string, msg *kafka.Message) {
	object.subscribeGroupMutex.RLock()
	defer object.subscribeGroupMutex.RUnlock()
	if group, ok := object.subscribeGroup[topic]; ok {
		for _, handler := range group {
			handler.OnMQMessage(msg.Value, msg.Offset)
		}
	}
}

// UseOption use option for MQSrv
func (object *MQSrv) UseOption(options ...MQSrvOption) *MQSrv {
	for _, option := range options {
		option(object)
	}
	return object
}

// Run run service
func (object *MQSrv) Run(ctx context.Context,
	wg *sync.WaitGroup,
	errCh chan<- error) {
	err := object.dial()
	errCh <- err
	if nil != err {
		return
	}
	object.kafkaConnMapMutex.RLock()
	defer object.kafkaConnMapMutex.RUnlock()
	for topic, conn := range object.kafkaConnMap {
		wg.Add(1)
		go func(topic string, conn *kafka.Conn) {
			defer wg.Done()
			log.Inst().Info("read kafka mq", zap.String("topic", topic))
			var msg kafka.Message
		loop:
			for {
				select {
				case <-ctx.Done():
					break loop
				default:
				}
				msg, err = object.Read(topic, 0, -1,
					defaultReadMaxBytes,
					time.Now().Add(defaultReadTimeout))
				if nil != err {
					if e, ok := err.(kafka.Error); ok && e.Timeout() {
						continue
					}
					log.Inst().Error("kafka read error", zap.String("error", err.Error()))
					break loop
				}
				object.Publish(topic, &msg)
			}
			log.Inst().Info("mq srv read loop done", zap.String("topic", topic))
		}(topic, conn)
	}
}

// Stop stop service
func (object *MQSrv) Stop() {
	object.BaseSrv.Stop()
	object.clean()
}

// SubscribeMessageHandler subscribe topic
func (object *MQSrv) SubscribeMessageHandler(topic string, handler types.IMQMessageHandler) *MQSrv {
	object.subscribeGroupMutex.Lock()
	defer object.subscribeGroupMutex.Unlock()
	object.subscribeGroup[topic] = append(object.subscribeGroup[topic], handler)
	return object
}

// SubscribeMessageCallback subscribe topic
func (object *MQSrv) SubscribeMessageCallback(topic string, callback types.MQMessageCallback) *MQSrv {
	return object.SubscribeMessageHandler(topic, mqMessageHandler.NewMQMessageHandler(callback))
}

// Write write raw to kafka
func (object *MQSrv) Write(topic string, raw []byte) (n int, err error) {
	object.kafkaConnMapMutex.RLock()
	defer object.kafkaConnMapMutex.RUnlock()
	conn, ok := object.kafkaConnMap[topic]
	if !ok {
		err = types.MQSrvTopicNotExistsError
		return
	}
	if err = conn.SetWriteDeadline(time.Now().Add(defaultReadTimeout)); nil != err {
		return
	}
	n, err = conn.Write(raw)
	return
}

// Read read message from kafka
func (object *MQSrv) Read(topic string,
	offset int64, whence int,
	maxBytes int,
	deadline time.Time) (msg kafka.Message, err error) {
	object.kafkaConnMapMutex.RLock()
	defer object.kafkaConnMapMutex.RUnlock()
	conn, ok := object.kafkaConnMap[topic]
	if !ok {
		err = types.MQSrvTopicNotExistsError
		return
	}
	switch whence {
	case kafka.SeekStart, kafka.SeekAbsolute, kafka.SeekEnd, kafka.SeekCurrent:
		if _, err = conn.Seek(offset, whence); nil != err {
			return
		}
	}
	if err = conn.SetReadDeadline(deadline); nil != err {
		return
	}
	msg, err = conn.ReadMessage(maxBytes)
	return
}

// MQSrvInst singleton
func MQSrvInst(options ...MQSrvOption) *MQSrv {
	mqSrvOnce.Do(func() {
		mqSrvInst = NewMQSrv(options...)
	})
	return mqSrvInst
}

// init initialize method
func init() {
	srvGroup.Inst().AddSrv(MQSrvInst())
}
