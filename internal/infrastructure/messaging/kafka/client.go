package kafka

import (
	"accounter/config"
	"accounter/internal/domain/event"
	"accounter/pkg/logger"
	"context"
	"fmt"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

// KafkaBroker Kafka queue broker
type KafkaBroker struct {
	ctx        context.Context
	logger     logger.Logger
	autoCommit bool

	reader *kafka.Reader
	writer *kafka.Writer

	readTimeout  time.Duration
	writeTimeout time.Duration

	subscribers []event.EventSubscriber
}

// NewBroker creates new KafkaBroker
func NewBroker(ctx context.Context, config config.Config, logger logger.Logger) *KafkaBroker {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     config.Kafka.Brokers,
		Topic:       config.Kafka.Topic,
		MaxBytes:    10e6,
		GroupID:     config.Kafka.Group,
		StartOffset: kafka.LastOffset,
	})

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   config.Kafka.Topic,
	})

	return &KafkaBroker{
		ctx:          ctx,
		reader:       reader,
		writer:       writer,
		readTimeout:  config.Kafka.ReadTimeout,
		writeTimeout: config.Kafka.WriteTimeout,
		autoCommit:   config.Kafka.AutoCommit,
		logger:       logger,
		subscribers:  make([]event.EventSubscriber, 0),
	}
}

// Name returns Kafka broker information
func (r *KafkaBroker) Name() string {
	return fmt.Sprintf("Kafka MQ broker on brokers: %v", r.reader.Config().Brokers)
}

// RegisterSubscribers register event subscribers
func (r *KafkaBroker) RegisterSubscribers(subscribers ...event.EventSubscriber) {
	r.subscribers = append(r.subscribers, subscribers...)
}

// Close reader/writer streams
func (r *KafkaBroker) Close() error {
	r.writer.Close()

	return r.reader.Close()
}
