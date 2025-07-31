package rocketevents

import (
	"encoding/json"
	"fmt"

	"github.com/soulcodex/rockets-message-processor/pkg/messaging"
)

const (
	RocketExplodedType = "RocketExploded"
)

type RocketExploded struct {
	*messaging.BaseMessage

	reason string
}

func (e *RocketExploded) Schema() string {
	return "rocket_launched"
}

func (e *RocketExploded) Reason() string {
	return e.reason
}

func (e *RocketExploded) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package
	if string(data) == rocketMessageNullContent || string(data) == `""` {
		return nil
	}

	type rocketLaunchedData struct {
		Reason string `json:"reason"`
	}

	type rocketLaunchedAlias struct {
		Metadata json.RawMessage    `json:"metadata"`
		Message  rocketLaunchedData `json:"message"`
	}

	alias := rocketLaunchedAlias{}
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("rocket exploded unmarshalling failed: %w", err)
	}

	var metadata RocketMessageMetadata
	err := json.Unmarshal(alias.Metadata, &metadata)
	if err != nil {
		return fmt.Errorf("rocket exploded metadata conversion failed: %w", err)
	}

	baseMessage, err := messaging.NewBaseMessage(
		RocketMessageID(metadata),
		RocketExplodedType,
		alias.Metadata,
		data,
		metadata.MessageTime,
	)
	if err != nil {
		return fmt.Errorf("rocket exploded base message creation failed: %w", err)
	}

	e.BaseMessage = baseMessage
	e.reason = alias.Message.Reason

	return nil
}
