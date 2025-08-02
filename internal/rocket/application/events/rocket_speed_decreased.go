package rocketevents

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	RocketSpeedDecreasedType = "RocketSpeedDecreased"
)

type RocketSpeedDecreased struct {
	EventID    string
	RocketID   string
	Amount     float64
	OccurredOn time.Time
}

func (e *RocketSpeedDecreased) Identifier() string {
	return e.EventID
}

func (e *RocketSpeedDecreased) BlockingKey() string {
	return "rocket" + ":" + e.RocketID
}

func (e *RocketSpeedDecreased) Type() string {
	return RocketSpeedDecreasedType
}

func (e *RocketSpeedDecreased) FromRawEvent(rm *RocketEventRaw) error {
	type rocketSpeedDecreasedData struct {
		Amount float64 `json:"by"`
	}

	var content rocketSpeedDecreasedData
	err := json.Unmarshal(rm.Message, &content)
	if err != nil {
		return fmt.Errorf("rocket speed decreased content conversion failed: %w", err)
	}

	e.EventID = rm.EventID()
	e.RocketID = rm.Metadata.Channel
	e.Amount = content.Amount
	e.OccurredOn = rm.Metadata.MessageTime

	return nil
}
