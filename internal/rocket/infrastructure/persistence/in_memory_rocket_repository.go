package rocketpersistence

import (
	"context"
	"slices"
	"sync"

	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
	"github.com/soulcodex/rockets-message-processor/pkg/errutil"
)

var (
	errRocketCannotBeNil = errutil.NewError("rocket cannot be nil")
)

type InMemoryRocketRepository struct {
	mutex   sync.RWMutex
	rockets map[rocketdomain.RocketID]*rocketdomain.Rocket
}

func NewInMemoryRocketRepository() *InMemoryRocketRepository {
	return &InMemoryRocketRepository{
		rockets: make(map[rocketdomain.RocketID]*rocketdomain.Rocket),
		mutex:   sync.RWMutex{},
	}
}

func (r *InMemoryRocketRepository) Find(_ context.Context, id rocketdomain.RocketID) (*rocketdomain.Rocket, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	rocket, exists := r.rockets[id]
	if !exists || rocket.Primitives().DeletedAt != nil {
		return nil, rocketdomain.NewRocketNotFoundError(id)
	}

	return rocket, nil
}

func (r *InMemoryRocketRepository) Search(_ context.Context, sortBy string, asc bool) (rocketdomain.RocketCollection, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(r.rockets) == 0 {
		return rocketdomain.NewRocketCollection(), nil
	}

	rockets := make([]*rocketdomain.Rocket, 0, len(r.rockets))
	for _, rocket := range r.rockets {
		if rocket.Primitives().DeletedAt == nil {
			rockets = append(rockets, rocket)
		}
	}

	slices.SortFunc(rockets, r.sortFunc(sortBy, asc))

	return rocketdomain.NewRocketCollection(rockets...), nil
}

func (r *InMemoryRocketRepository) Save(_ context.Context, rocket *rocketdomain.Rocket) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if rocket == nil {
		return rocketdomain.NewRocketStoreError().Wrap(errRocketCannotBeNil)
	}

	r.rockets[rocket.ID()] = rocket
	return nil
}

func (r *InMemoryRocketRepository) sortFunc(sortBy string, asc bool) func(one, two *rocketdomain.Rocket) int {
	switch sortBy {
	case "created_at", "updated_at":
		return r.sortByTimestamps(sortBy, asc)
	case "launch_speed":
		return r.sortByLaunchSpeed(asc)
	default:
		return func(_, _ *rocketdomain.Rocket) int {
			return 1
		}
	}
}

func (r *InMemoryRocketRepository) sortByTimestamps(sortBy string, asc bool) func(left, right *rocketdomain.Rocket) int {
	return func(left, right *rocketdomain.Rocket) int {
		switch sortBy {
		case "created_at":
			if asc {
				return left.Primitives().CreatedAt.Compare(right.Primitives().CreatedAt)
			}
			return right.Primitives().CreatedAt.Compare(left.Primitives().CreatedAt)
		case "updated_at":
			if asc {
				return left.Primitives().UpdatedAt.Compare(right.Primitives().UpdatedAt)
			}

			return right.Primitives().UpdatedAt.Compare(left.Primitives().UpdatedAt)
		default:
			return 1 // Default case if sortBy is not recognized
		}
	}
}

func (r *InMemoryRocketRepository) sortByLaunchSpeed(asc bool) func(left, right *rocketdomain.Rocket) int {
	return func(left, right *rocketdomain.Rocket) int {
		leftSpeed, rightSpeed := left.Primitives().LaunchSpeed, right.Primitives().LaunchSpeed

		if leftSpeed == rightSpeed {
			return 0
		}

		if asc {
			if leftSpeed < rightSpeed {
				return -1
			}
			return 1
		}

		if leftSpeed > rightSpeed {
			return -1
		}
		return 1
	}
}
