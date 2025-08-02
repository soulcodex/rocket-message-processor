package rocketevents

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	RocketLaunchedType = "RocketLaunched"
)

type RocketLaunched struct {
	EventID     string
	RocketID    string
	RocketType  string
	LaunchSpeed int64
	Mission     string
	OccurredOn  time.Time
}

func (e *RocketLaunched) Identifier() string {
	return e.EventID
}

func (e *RocketLaunched) BlockingKey() string {
	return "rocket" + ":" + e.RocketID
}

func (e *RocketLaunched) Type() string {
	return RocketLaunchedType
}

func (e *RocketLaunched) FromRawEvent(rm *RocketEventRaw) error {
	type rocketLaunchedData struct {
		RocketType  string `json:"type"`
		LaunchSpeed int64  `json:"launchSpeed"`
		Mission     string `json:"mission"`
	}

	var content rocketLaunchedData
	err := json.Unmarshal(rm.Message, &content)
	if err != nil {
		return fmt.Errorf("rocket launched content conversion failed: %w", err)
	}

	e.EventID = rm.EventID()
	e.RocketID = rm.Metadata.Channel
	e.RocketType = content.RocketType
	e.LaunchSpeed = content.LaunchSpeed
	e.Mission = content.Mission
	e.OccurredOn = rm.Metadata.MessageTime

	return nil
}
