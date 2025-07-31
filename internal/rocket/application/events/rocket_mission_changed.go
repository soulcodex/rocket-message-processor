package rocketevents

import (
	"encoding/json"
	"fmt"

	"github.com/soulcodex/rockets-message-processor/pkg/messaging"
)

const (
	RocketMissionChangedType = "RocketMissionChanged"
)

type RocketMissionChanged struct {
	*messaging.BaseMessage

	newMission string
}

func (e *RocketMissionChanged) Schema() string {
	return "rocket_mission_changed"
}

func (e *RocketMissionChanged) NewMission() string {
	return e.newMission
}

func (e *RocketMissionChanged) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	type rocketSpeedIncreasedData struct {
		Mission string `json:"newMission"`
	}

	type rocketLaunchedAlias struct {
		Metadata json.RawMessage          `json:"metadata"`
		Message  rocketSpeedIncreasedData `json:"message"`
	}

	alias := rocketLaunchedAlias{}
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("rocker mission changed unmarshalling failed: %w", err)
	}

	var metadata RocketMessageMetadata
	err := json.Unmarshal(alias.Metadata, &metadata)
	if err != nil {
		return fmt.Errorf("rocker mission changed metadata conversion failed: %w", err)
	}

	baseMessage, err := messaging.NewBaseMessage(
		RocketMessageID(metadata),
		RocketMissionChangedType,
		alias.Metadata,
		data,
		metadata.MessageTime,
	)
	if err != nil {
		return fmt.Errorf("rocker mission changed base message creation failed: %w", err)
	}

	e.BaseMessage = baseMessage
	e.newMission = alias.Message.Mission

	return nil
}
