package rocketentrypoint

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	rocketevents "github.com/soulcodex/rockets-message-processor/internal/rocket/application/events"
	"github.com/soulcodex/rockets-message-processor/pkg/bus"
	eventbus "github.com/soulcodex/rockets-message-processor/pkg/bus/event"
	distributedsync "github.com/soulcodex/rockets-message-processor/pkg/distributed-sync"
	"github.com/soulcodex/rockets-message-processor/pkg/errutil"
	httpserver "github.com/soulcodex/rockets-message-processor/pkg/http-server"
	"github.com/soulcodex/rockets-message-processor/pkg/messaging"
)

func HandleReceiveRocketMessageV1HTTP(
	eventBus eventbus.EventBus,
	mutex distributedsync.MutexService,
	deduplicator messaging.Deduplicator,
	responseWriter *httpserver.JSONResponseWriter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			responseWriter.WriteErrorResponse(
				r.Context(),
				w,
				[]string{"failed to read request body"},
				http.StatusBadRequest,
			)
			return
		}

		var raw rocketevents.RocketEventRaw
		unmarshalErr := json.Unmarshal(body, &raw)
		if unmarshalErr != nil {
			responseWriter.WriteErrorResponse(
				r.Context(),
				w,
				[]string{"failed to parse request body"},
				http.StatusBadRequest,
			)
			return
		}

		rocketEvent, err := rocketevents.ResolveRocketEvent(raw.Metadata.MessageType, &raw)
		if err != nil {
			responseWriter.WriteErrorResponse(
				r.Context(),
				w,
				[]string{"unable to resolve rocket event type"},
				http.StatusBadRequest,
			)
			return
		}

		if handleErr := handleEventWithDeduplication(
			r.Context(),
			eventBus,
			mutex,
			deduplicator,
			rocketEvent,
		); handleErr != nil {
			responseWriter.WriteErrorResponse(
				r.Context(),
				w,
				[]string{handleErr.Error()},
				http.StatusInternalServerError,
			)
			return
		}

		responseWriter.WriteResponse(r.Context(), w, nil, http.StatusNoContent)
	}
}

func handleEventWithDeduplication(
	ctx context.Context,
	eventBus eventbus.EventBus,
	mutex distributedsync.MutexService,
	deduplicator messaging.Deduplicator,
	rocketEvent rocketevents.RocketEvent,
) error {
	if deduplicator == nil {
		return errutil.NewError("deduplicator is not configured")
	}

	if err := checkIfRocketEventIsDuplicated(ctx, deduplicator, rocketEvent); err != nil {
		return fmt.Errorf("duplicated check failed: %w", err)
	}

	if blockingDto, match := rocketEvent.(bus.BlockingDto); match {
		if err := bus.DispatchBlocking(eventBus, mutex)(ctx, blockingDto); err != nil {
			return fmt.Errorf("failed to dispatch blocking event: %w", err)
		}
	}

	if err := markRocketEventAsProcessed(ctx, deduplicator, rocketEvent); err != nil {
		return fmt.Errorf("failed to mark event as processed: %w", err)
	}

	return nil
}

func checkIfRocketEventIsDuplicated(
	ctx context.Context,
	deduplicator messaging.Deduplicator,
	rocketEvent rocketevents.RocketEvent,
) error {
	if deduplicator == nil {
		return errutil.NewError("deduplicator is not configured")
	}

	message, match := rocketEvent.(messaging.Message)
	if !match {
		return errutil.NewError("rocket event is not a message")
	}

	if duplicated, err := deduplicator.IsDuplicate(ctx, message); err == nil && !duplicated {
		return nil
	}

	return fmt.Errorf("rocket event is duplicated: %s", message.Identifier())
}

func markRocketEventAsProcessed(
	ctx context.Context,
	deduplicator messaging.Deduplicator,
	rocketEvent rocketevents.RocketEvent,
) error {
	if deduplicator == nil {
		return errutil.NewError("deduplicator is not configured")
	}

	message, match := rocketEvent.(messaging.Message)
	if !match {
		return errutil.NewError("rocket event is not a message")
	}

	markErr := deduplicator.MarkProcessed(ctx, message)
	if markErr != nil {
		return fmt.Errorf("failed to mark rocket event as processed: %w", markErr)
	}

	return nil
}
