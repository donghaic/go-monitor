package common

import "errors"

var (
	Redis_Init_Error = errors.New("redis init error")

	Redis_PubSub_Error = errors.New("redis pubsub error")

	Kafka_Init_Error = errors.New("kafka init error")
)
