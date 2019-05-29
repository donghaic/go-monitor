package kafka

import (
	"encoding/json"
	"errors"
	"fmt"
	"gt-monitor/common/zap"
	"github.com/Shopify/sarama"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	ON  int32 = 1
	OFF int32 = 0
)

type Producer interface {
	SendKeyedMessage(topic, key string, value interface{}) (err error)

	SendMessage(topic string, value interface{}) (err error)

	GetConf() *Options
}

type client struct {
	producer sarama.AsyncProducer
}

func (c *client) close() error {
	return c.producer.Close()
}

type kafkaAsyncProducer struct {
	client *client
	conf   *Options
	status int32
}

func (k kafkaAsyncProducer) SendKeyedMessage(topic, key string, value interface{}) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(data),
		Timestamp: time.Now(),
	}

	k.client.producer.Input() <- msg

	return nil

}
func (k kafkaAsyncProducer) SendMessage(topic string, value interface{}) (err error) {

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if atomic.LoadInt32(&k.status) == OFF {
		zap.Get().Error("kafka producer on OFF status, data=", data)
		return errors.New("product is on OFF status")
	}

	message := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(data),
		Timestamp: time.Now(),
	}

	k.client.producer.Input() <- message

	return nil

}

func (k kafkaAsyncProducer) GetConf() *Options {
	return k.conf
}

func (k kafkaAsyncProducer) handleRebuild() {
	// return.Successes,return.Errors 为true 需要处理error 和 succedss,否则会阻塞
	go func() {
		for {
			select {
			case <-k.client.producer.Successes():
				// discard
			case err := <-k.client.producer.Errors():
				if err != nil {
					fmt.Println("Send message to kafka error. ", err)
					zap.Get().Error("send to kafka error", err)

					// 使用cas设置标识位OFF， 重建完成后再设成ON
					if atomic.CompareAndSwapInt32(&k.status, ON, OFF) {
						// 重建
						zap.Get().Info("bad things happened, start recreate kafka producer")
						go func() {
							newClient, err := reCreateProducer(k.conf)
							if err == nil {
								oldClient := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&k.client)), unsafe.Pointer(newClient))
								atomic.StoreInt32(&k.status, ON)
								(*client)(oldClient).close()
							}
						}()
					}
				}
			}
		}
	}()
}

func
reCreateProducer(config *Options) (*client, error) {
	zap.Get().Info("rebuild kafka client ...")
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Retry.Max = 10
	cfg.Producer.Timeout = 1000
	cfg.ClientID = config.ClientID

	saramaProducer, err := sarama.NewAsyncProducer(config.Brokers, cfg)

	if err == nil {
		return &client{producer: saramaProducer}, err
	} else {
		reCreateProducer(config)
	}
	return nil, err

}

// New ...
func New(config *Options) (Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Compression = sarama.CompressionGZIP     // Compress messages
	cfg.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Retry.Max = 10
	cfg.ClientID = config.ClientID

	client, err := reCreateProducer(config)

	p := new(kafkaAsyncProducer)
	p.conf = config
	p.client = client
	atomic.StoreInt32(&p.status, ON)
	p.handleRebuild()

	return p, err
}
