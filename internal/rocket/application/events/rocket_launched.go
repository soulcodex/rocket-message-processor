package rocketevents

import (
	"encoding/json"
	"fmt"

	"github.com/soulcodex/rockets-message-processor/pkg/messaging"
)

const (
	RocketLaunchedType = "RocketLaunched"
)

type RocketLaunched struct {
	*messaging.BaseMessage

	rocketType  string
	launchSpeed uint64
	mission     string
}

func (e *RocketLaunched) Schema() string {
	return "rocket_launched"
}

func (e *RocketLaunched) RocketType() string {
	return e.rocketType
}

func (e *RocketLaunched) LaunchSpeed() uint64 {
	return e.launchSpeed
}

func (e *RocketLaunched) Mission() string {
	return e.mission
}

func (e *RocketLaunched) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package
	if string(data) == rocketMessageNullContent || string(data) == `""` {
		return nil
	}

	type rocketLaunchedData struct {
		RocketType  string `json:"rocketType"`
		LaunchSpeed uint64 `json:"launchSpeed"`
		Mission     string `json:"mission"`
	}

	type rocketLaunchedAlias struct {
		Metadata json.RawMessage    `json:"metadata"`
		Message  rocketLaunchedData `json:"message"`
	}

	alias := rocketLaunchedAlias{}
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("rocket launched unmarshalling failed: %w", err)
	}

	var metadata RocketMessageMetadata
	err := json.Unmarshal(alias.Metadata, &metadata)
	if err != nil {
		return fmt.Errorf("rocket launched metadata conversion failed: %w", err)
	}

	baseMessage, err := messaging.NewBaseMessage(
		RocketMessageID(metadata),
		RocketLaunchedType,
		alias.Metadata,
		data,
		metadata.MessageTime,
	)
	if err != nil {
		return fmt.Errorf("rocket launched base message creation failed: %w", err)
	}

	e.BaseMessage = baseMessage
	e.rocketType = alias.Message.RocketType
	e.launchSpeed = alias.Message.LaunchSpeed
	e.mission = alias.Message.Mission

	return nil
}
