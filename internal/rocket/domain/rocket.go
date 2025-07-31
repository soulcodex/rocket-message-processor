package rocketdomain

import (
	"time"
)

type Rocket struct {
	id          RocketID
	rocketType  RocketType
	launchSpeed LaunchSpeed
	mission     Mission
	createdAt   time.Time
	updatedAt   time.Time
}

func NewRocket(
	id RocketID,
	rocketType RocketType,
	launchSpeed LaunchSpeed,
	mission Mission,
	at time.Time,
) *Rocket {
	return &Rocket{
		id:          id,
		rocketType:  rocketType,
		launchSpeed: launchSpeed,
		mission:     mission,
		createdAt:   at,
		updatedAt:   at,
	}
}

func (r *Rocket) Primitives() RocketPrimitives {
	return FromDomain(r)
}

func (r *Rocket) IncreaseLaunchSpeed(increase LaunchSpeed, at time.Time) {
	if at.IsZero() || at.After(r.updatedAt) {
		return
	}

	r.launchSpeed += increase
	r.updatedAt = at
}

func (r *Rocket) ChangeMission(newMission Mission, at time.Time) {
	if at.IsZero() || at.After(r.updatedAt) {
		return
	}

	r.mission = newMission
	r.updatedAt = at
}
