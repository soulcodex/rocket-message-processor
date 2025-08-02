package rocketevents

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/soulcodex/rockets-message-processor/pkg/bus"
)

const (
	rocketMessageNullContent = "null"
)

type RocketEvent interface {
	bus.Dto

	FromRawEvent(rm *RocketEventRaw) error
}

type RocketEventMetadata struct {
	Channel       string    `json:"channel"`
	MessageNumber uint64    `json:"messageNumber"`
	MessageTime   time.Time `json:"messageTime"`
	MessageType   string    `json:"messageType"`
}

func (m *RocketEventMetadata) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package
	if string(data) == rocketMessageNullContent || string(data) == `""` {
		return nil
	}

	type metadataAlias struct {
		Channel       string `json:"channel"`
		MessageNumber uint64 `json:"messageNumber"`
		MessageTime   string `json:"messageTime"`
		MessageType   string `json:"messageType"`
	}

	var metadata metadataAlias
	if err := json.Unmarshal(data, &metadata); err != nil {
		return fmt.Errorf("rocket message metadata unmarshalling failed: %w", err)
	}

	var messageTime time.Time
	if metadata.MessageTime != "" {
		timeParsed, timeParseErr := time.Parse(time.RFC3339, metadata.MessageTime)
		if timeParseErr != nil {
			return fmt.Errorf("rocket message metadata time parsing failed: %w", timeParseErr)
		}
		messageTime = timeParsed
	}

	m.Channel = metadata.Channel
	m.MessageNumber = metadata.MessageNumber
	m.MessageTime = messageTime
	m.MessageType = metadata.MessageType

	return nil
}

type RocketEventRaw struct {
	Metadata RocketEventMetadata `json:"metadata"`
	Message  json.RawMessage     `json:"message"`
}

// EventID generates a unique identifier for the event based on its metadata.
func (rm *RocketEventRaw) EventID() string {
	return fmt.Sprintf(
		"%s:%s:%d:%d",
		rm.Metadata.Channel,
		strings.ToLower(rm.Metadata.MessageType),
		rm.Metadata.MessageNumber,
		rm.Metadata.MessageTime.UnixMilli(),
	)
}
