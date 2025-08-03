package rockettest

import (
	"testing"
	"time"

	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
)

type RocketMotherOpt func(*RocketMother)

func WithRocketID(id string) RocketMotherOpt {
	return func(m *RocketMother) {
		m.primitives.ID = id
	}
}

func WithLaunchSpeed(speed int64) RocketMotherOpt {
	return func(m *RocketMother) {
		m.primitives.LaunchSpeed = speed
	}
}

func WithSoftDeletion() RocketMotherOpt {
	return func(m *RocketMother) {
		now := time.Now()
		m.primitives.DeletedAt = &now
	}
}

func WithUpdateDate(at time.Time) RocketMotherOpt {
	return func(m *RocketMother) {
		m.primitives.UpdatedAt = at
	}
}

type RocketMother struct {
	primitives rocketdomain.RocketPrimitives
}

func NewRocketMother(opts ...RocketMotherOpt) *RocketMother {
	mother := &RocketMother{
		primitives: newRocketPrimitives(),
	}

	for _, opt := range opts {
		opt(mother)
	}

	return mother
}

func (m *RocketMother) Build(t *testing.T) *rocketdomain.Rocket {
	t.Helper()

	at := m.primitives.CreatedAt
	if m.primitives.UpdatedAt.After(at) {
		at = m.primitives.UpdatedAt
	}

	rocket := rocketdomain.NewRocket(
		rocketdomain.RocketID(m.primitives.ID),
		rocketdomain.RocketType(m.primitives.RocketType),
		rocketdomain.LaunchSpeed(m.primitives.LaunchSpeed),
		rocketdomain.Mission(m.primitives.Mission),
		at,
	)

	if m.primitives.DeletedAt != nil {
		rocket.Delete(time.Now())
	}

	return rocket
}

func newRocketPrimitives() rocketdomain.RocketPrimitives {
	const defaultLaunchSpeed = 5000
	return rocketdomain.RocketPrimitives{
		ID:          "122c31b0-a3c4-411a-bc07-5f342f0d78e4",
		RocketType:  "Falcon 9",
		LaunchSpeed: defaultLaunchSpeed,
		Mission:     "ARTEMIS",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
}
