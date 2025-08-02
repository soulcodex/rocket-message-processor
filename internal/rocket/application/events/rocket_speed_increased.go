package rocketevents

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	RocketSpeedIncreasedType = "RocketSpeedIncreased"
)

type RocketSpeedIncreased struct {
	EventID    string
	RocketID   string
	Amount     float64
	OccurredOn time.Time
}

func (e *RocketSpeedIncreased) Identifier() string {
	return e.EventID
}

func (e *RocketSpeedIncreased) BlockingKey() string {
	return "rocket" + ":" + e.RocketID
}

func (e *RocketSpeedIncreased) Type() string {
	return RocketSpeedIncreasedType
}

func (e *RocketSpeedIncreased) FromRawEvent(rm *RocketEventRaw) error {
	type rocketSpeedIncreasedData struct {
		Amount float64 `json:"by"`
	}

	var content rocketSpeedIncreasedData
	err := json.Unmarshal(rm.Message, &content)
	if err != nil {
		return fmt.Errorf("rocket speed increased content conversion failed: %w", err)
	}

	e.EventID = rm.EventID()
	e.RocketID = rm.Metadata.Channel
	e.Amount = content.Amount
	e.OccurredOn = rm.Metadata.MessageTime

	return nil
}
