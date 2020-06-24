package types

import "errors"

var (
	//ErrMQSrvTopicNotExists mq topic not exists error
	ErrMQSrvTopicNotExists = errors.New("topic not exists")
)
