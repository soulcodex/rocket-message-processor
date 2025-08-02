package rocketevents

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	RocketMissionChangedType = "RocketMissionChanged"
)

type RocketMissionChanged struct {
	EventID    string
	RocketID   string
	NewMission string
	OccurredOn time.Time
}

func (e *RocketMissionChanged) Identifier() string {
	return e.EventID
}

func (e *RocketMissionChanged) BlockingKey() string {
	return "rocket" + ":" + e.RocketID
}

func (e *RocketMissionChanged) Type() string {
	return RocketMissionChangedType
}

func (e *RocketMissionChanged) FromRawEvent(rm *RocketEventRaw) error {
	type rocketMissionChangedData struct {
		NewMission string `json:"newMission"`
	}

	var content rocketMissionChangedData
	err := json.Unmarshal(rm.Message, &content)
	if err != nil {
		return fmt.Errorf("rocket mission changed content conversion failed: %w", err)
	}

	e.EventID = rm.EventID()
	e.RocketID = rm.Metadata.Channel
	e.NewMission = content.NewMission
	e.OccurredOn = rm.Metadata.MessageTime

	return nil
}
