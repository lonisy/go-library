package kafka

import (
    kafka "github.com/segmentio/kafka-go"
    "strings"
    "time"
)

type KafkaStruct struct {
    Writer   *kafka.Writer
    WriterOn bool
}

var KafkaGo KafkaStruct

type KafkaConfig struct {
    Brokers   string `json:"brokers"`
    GroupID   string `json:"group_id"`
    Topic     string `json:"topic"`
    Partition int    `json:"partition"`
    MinBytes  int    `json:"min_bytes"`
    MaxBytes  int    `json:"max_bytes"`
}

func (r *KafkaStruct) Consumer(config KafkaConfig) *kafka.Reader {
    brokers := strings.Split(config.Brokers, ",")
    return kafka.NewReader(kafka.ReaderConfig{
        Brokers:        brokers,
        GroupID:        config.GroupID,
        Topic:          config.Topic,
        Partition:      config.Partition,
        CommitInterval: 1 * time.Second,
        MaxBytes:       10e6,
        QueueCapacity:  5000,
    })
}

func (r *KafkaStruct) Producer(config KafkaConfig) *kafka.Writer {
    brokers := strings.Split(config.Brokers, ",")
    r.Writer = &kafka.Writer{
        Addr:         kafka.TCP(brokers...),
        Topic:        config.Topic,
        Balancer:     &kafka.LeastBytes{},
        BatchTimeout: 100 * time.Millisecond,
        BatchSize:    2000,
        RequiredAcks: kafka.RequireOne,
    }
    return r.Writer
}
