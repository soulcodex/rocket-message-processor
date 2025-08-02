package rocketevents

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	RocketExplodedType = "RocketExploded"
)

type RocketExploded struct {
	EventID    string
	RocketID   string
	Reason     string
	OccurredOn time.Time
}

func (e *RocketExploded) Identifier() string {
	return e.EventID
}

func (e *RocketExploded) BlockingKey() string {
	return "rocket" + ":" + e.RocketID
}

func (e *RocketExploded) Type() string {
	return RocketExplodedType
}

func (e *RocketExploded) FromRawEvent(rm *RocketEventRaw) error {
	type rocketExplodedData struct {
		Reason string `json:"reason"`
	}

	var content rocketExplodedData
	err := json.Unmarshal(rm.Message, &content)
	if err != nil {
		return fmt.Errorf("rocket exploded content conversion failed: %w", err)
	}

	e.EventID = rm.EventID()
	e.RocketID = rm.Metadata.Channel
	e.Reason = content.Reason
	e.OccurredOn = rm.Metadata.MessageTime

	return nil
}
