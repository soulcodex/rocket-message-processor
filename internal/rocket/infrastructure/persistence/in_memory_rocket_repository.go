package rocketpersistence

import (
	"context"
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
	if !exists {
		return nil, rocketdomain.NewRocketNotFoundError(id)
	}

	return rocket, nil
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
