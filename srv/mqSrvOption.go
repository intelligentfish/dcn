package srv

// MQSrvOption
type MQSrvOption func(object *MQSrv)

// KafkaParamOption kafka param option
func KafkaParamOption(param *KafkaParam) MQSrvOption {
	return func(object *MQSrv) {
		object.kafkaParamList = append(object.kafkaParamList, param)
	}
}
