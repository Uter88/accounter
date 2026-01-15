package kafka

import (
	"accounter/internal/domain/event"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Run start to read Event queuies
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

			err = r.sendMessageToSubscribers(msg)

			if err == nil || r.autoCommit {
				r.reader.CommitMessages(ctx, msg)
			} else {
				r.logger.Errorf("error process message: %s", err)
			}
		}
	}
}

// sendMessageToSubscribers decode Message to Event and try to send to all Event subscribers
func (r *KafkaBroker) sendMessageToSubscribers(msg kafka.Message) error {
	var event event.Event

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("fail to unmarshall message to Event: %w", err)
	}

	var errs []error

	for _, s := range r.subscribers {
		if err := s.SubscribeEvent(event); err != nil {
			errs = append(errs, fmt.Errorf("-%s:%w-", s.Name(), err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("fail to send Event to subscriber(s): %s", errors.Join(errs...))
	}

	return nil
}
