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
	deletedAt   *time.Time
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
	return primitivesFromDomain(r)
}

func (r *Rocket) ID() RocketID {
	return r.id
}

func (r *Rocket) ChangeLaunchSpeed(speed LaunchSpeed, at time.Time) {
	if r.deletedAt != nil || at.IsZero() || at.After(r.updatedAt) {
		return
	}

	if speed < 0 && r.launchSpeed+speed < 0 {
		r.launchSpeed = 0
		r.updatedAt = at
		return
	}

	r.launchSpeed += speed
	r.updatedAt = at
}

func (r *Rocket) ChangeMission(newMission Mission, at time.Time) {
	if r.deletedAt != nil || at.IsZero() || at.After(r.updatedAt) {
		return
	}

	r.mission = newMission
	r.updatedAt = at
}

func (r *Rocket) Delete(at time.Time) {
	if r.deletedAt != nil || at.IsZero() || at.After(r.updatedAt) {
		return
	}

	r.updatedAt = at
	r.deletedAt = &at
}
