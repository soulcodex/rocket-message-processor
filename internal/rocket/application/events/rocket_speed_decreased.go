package rocketevents

import (
	"encoding/json"
	"fmt"

	"github.com/soulcodex/rockets-message-processor/pkg/messaging"
)

const (
	RocketSpeedDecreasedType = "RocketSpeedDecreased"
)

type RocketSpeedDecreased struct {
	*messaging.BaseMessage

	amount float64
}

func (e *RocketSpeedDecreased) Schema() string {
	return "rocket_speed_decreased"
}

func (e *RocketSpeedDecreased) Amount() float64 {
	return e.amount
}

func (e *RocketSpeedDecreased) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	type rocketSpeedIncreasedData struct {
		Amount float64 `json:"by"`
	}

	type rocketLaunchedAlias struct {
		Metadata json.RawMessage          `json:"metadata"`
		Message  rocketSpeedIncreasedData `json:"message"`
	}

	alias := rocketLaunchedAlias{}
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("rocket speed decreased unmarshalling failed: %w", err)
	}

	var metadata RocketMessageMetadata
	err := json.Unmarshal(alias.Metadata, &metadata)
	if err != nil {
		return fmt.Errorf("rocket speed decreased metadata conversion failed: %w", err)
	}

	baseMessage, err := messaging.NewBaseMessage(
		RocketMessageID(metadata),
		RocketSpeedDecreasedType,
		alias.Metadata,
		data,
		metadata.MessageTime,
	)
	if err != nil {
		return fmt.Errorf("rocket speed decreased base message creation failed: %w", err)
	}

	e.BaseMessage = baseMessage
	e.amount = alias.Message.Amount

	return nil
}
