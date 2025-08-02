package di

import (
	"context"

	rocketevents "github.com/soulcodex/rockets-message-processor/internal/rocket/application/events"
	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
	rocketentrypoint "github.com/soulcodex/rockets-message-processor/internal/rocket/infrastructure/entrypoint"
	rocketpersistence "github.com/soulcodex/rockets-message-processor/internal/rocket/infrastructure/persistence"
	"github.com/soulcodex/rockets-message-processor/pkg/bus"
	httpserver "github.com/soulcodex/rockets-message-processor/pkg/http-server"
)

type RocketModule struct {
	Repository rocketdomain.RocketRepository
	Creator    *rocketdomain.RocketCreator
	Updater    *rocketdomain.RocketUpdater
}

func NewRocketModule(_ context.Context, common *CommonServices) *RocketModule {
	rocketRepo := rocketpersistence.NewInMemoryRocketRepository()
	creator := rocketdomain.NewRocketCreator(rocketRepo)
	updater := rocketdomain.NewRocketUpdater(rocketRepo)

	common.Router.Post(
		"/messages",
		rocketentrypoint.HandleReceiveRocketMessageV1HTTP(
			common.EventBus,
			common.Mutex,
			common.Deduplicator,
			httpserver.NewJSONResponseMiddleware(common.Logger),
		),
	)

	launchEvtHandler := rocketevents.NewCreateRocketOnRocketLaunched(creator)
	bus.MustRegister(common.EventBus, &rocketevents.RocketLaunched{}, launchEvtHandler)

	explodeEvtHandler := rocketevents.NewDeleteRocketOnRocketExploded(updater)
	bus.MustRegister(common.EventBus, &rocketevents.RocketExploded{}, explodeEvtHandler)

	paramsChangeEvtHandler := rocketevents.NewUpdateRocketOnRocketParamsChanged(updater)
	bus.MustRegister(common.EventBus, &rocketevents.RocketMissionChanged{}, paramsChangeEvtHandler)
	bus.MustRegister(common.EventBus, &rocketevents.RocketSpeedIncreased{}, paramsChangeEvtHandler)
	bus.MustRegister(common.EventBus, &rocketevents.RocketSpeedDecreased{}, paramsChangeEvtHandler)

	return &RocketModule{
		Repository: rocketRepo,
		Creator:    creator,
		Updater:    updater,
	}
}
