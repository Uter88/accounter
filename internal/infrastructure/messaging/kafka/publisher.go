package kafka

import (
	"accounter/internal/domain/event"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

// PublishEvent publish new event to queue
func (r *KafkaBroker) PublishEvent(event event.Event) error {
	return r.SendMessages(event)
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
