package kafka

import (
	"accounter/internal/domain/event"
	"context"
	"encoding/json"
	"errors"
)

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
