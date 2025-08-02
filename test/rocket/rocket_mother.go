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

	return rocketdomain.NewRocket(
		rocketdomain.RocketID(m.primitives.ID),
		rocketdomain.RocketType(m.primitives.RocketType),
		rocketdomain.LaunchSpeed(m.primitives.LaunchSpeed),
		rocketdomain.Mission(m.primitives.Mission),
		m.primitives.CreatedAt,
	)
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
