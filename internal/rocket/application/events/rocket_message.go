package rocketevents

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	rocketMessageNullContent = "null"
)

type RocketMessageMetadata struct {
	Channel       string    `json:"channel"`
	MessageNumber uint64    `json:"messageNumber"`
	MessageTime   time.Time `json:"messageTime"`
	MessageType   string    `json:"messageType"`
}

func (m *RocketMessageMetadata) UnmarshalJSON(data []byte) error {
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

	*m = RocketMessageMetadata{
		Channel:       metadata.Channel,
		MessageNumber: metadata.MessageNumber,
		MessageTime:   messageTime,
		MessageType:   metadata.MessageType,
	}

	return nil
}

func RocketMessageID(m RocketMessageMetadata) string {
	return fmt.Sprintf("%s:%s:%d", m.Channel, strings.ToLower(m.MessageType), m.MessageNumber)
}
