package queues

import (
	"accounter/config"
	"accounter/internal/domain/event"
	"accounter/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

// KafkaBroker Kafka queue broker
type KafkaBroker struct {
	ctx    context.Context
	logger logger.Logger

	reader *kafka.Reader
	writer *kafka.Writer

	readTimeout  time.Duration
	writeTimeout time.Duration

	subscribers []event.EventSubscriber
}

// NewBroker creates new KafkaBroker
func NewBroker(ctx context.Context, config config.Config, logger logger.Logger) *KafkaBroker {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  config.Kafka.Brokers,
		Topic:    config.Kafka.Topic,
		MaxBytes: 10e6,
		GroupID:  config.Kafka.Group,
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

// PublishEvent publish new event to queue
func (r *KafkaBroker) PublishEvent(event event.Event) error {
	return r.SendMessages(event)
}

// Run start event queue reader
func (r *KafkaBroker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			msg, err := r.reader.FetchMessage(r.ctx)

			if err != nil {
				r.logger.Errorf("error read event message: %s", err.Error())
				continue
			}

			var event event.Event

			if err = json.Unmarshal(msg.Value, &event); err != nil {
				r.logger.Errorf("error read event message: %s", err.Error())
				continue
			}

			var errs []error

			for _, subscriber := range r.subscribers {
				if err = subscriber.SubscribeEvent(event); err != nil {
					errs = append(errs, err)
				}
			}

			if len(errs) == 0 {
				r.reader.CommitMessages(ctx, msg)
			} else {
				r.logger.Errorf("error send event to subscriber(s): %s", errors.Join(errs...))
			}
		}
	}
}

// SendMessages send message
func (r *KafkaBroker) SendMessages(payload ...any) error {
	messages := make([][]byte, len(payload))

	for i := range payload {
		data, _ := json.Marshal(payload[i])
		messages[i] = data
	}

	return r.SendRawMessages(messages...)
}

// SendRawMessages send bytes messages to queue
func (r *KafkaBroker) SendRawMessages(payload ...[]byte) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.writeTimeout)
	defer cancel()

	messages := make([]kafka.Message, len(payload))

	for i := range payload {
		messages[i] = kafka.Message{Value: payload[i]}
	}

	return r.writer.WriteMessages(ctx, messages...)
}

// Close reader/writer streams
func (r *KafkaBroker) Close() error {
	r.writer.Close()

	return r.reader.Close()
}
