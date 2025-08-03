package di

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/soulcodex/rockets-message-processor/configs"
	eventbus "github.com/soulcodex/rockets-message-processor/pkg/bus/event"
	querybus "github.com/soulcodex/rockets-message-processor/pkg/bus/query"
	distributedsync "github.com/soulcodex/rockets-message-processor/pkg/distributed-sync"
	httpserver "github.com/soulcodex/rockets-message-processor/pkg/http-server"
	"github.com/soulcodex/rockets-message-processor/pkg/logger"
	"github.com/soulcodex/rockets-message-processor/pkg/messaging"
	"github.com/soulcodex/rockets-message-processor/pkg/utils"
)

type CommonServices struct {
	Config       *configs.Config
	Logger       logger.ZerologLogger
	RedisClient  *redis.Client
	EventBus     eventbus.Bus
	QueryBus     querybus.Bus
	Mutex        distributedsync.MutexService
	Deduplicator messaging.Deduplicator
	Router       *httpserver.Router
	UUIDProvider utils.UUIDProvider
	TimeProvider utils.DateTimeProvider
}

func MustInitCommonServices(ctx context.Context) *CommonServices {
	cfg, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	appLogger := logger.NewZerologLogger(
		ctx,
		"rocket-message-processor",
		logger.WithLogLevel(cfg.LogLevel),
		logger.WithAppVersion(cfg.AppVersion),
	)

	timeProvider := utils.NewSystemTimeProvider()

	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(redisOpts)

	routerOpts := []httpserver.RouterConfigFunc{
		httpserver.WithHost(cfg.HTTPHost),
		httpserver.WithPort(cfg.HTTPPort),
		httpserver.WithReadTimeoutSeconds(cfg.HTTPReadTimeout),
		httpserver.WithWriteTimeoutSeconds(cfg.HTTPWriteTimeout),
		httpserver.WithMiddleware(httpserver.NewPanicRecoverMiddleware(appLogger).Middleware),
		httpserver.WithMiddleware(httpserver.NewRequestLoggingMiddleware(appLogger, timeProvider).Middleware),
		httpserver.WithCORSMiddleware(),
	}
	router := httpserver.New(routerOpts...)

	eventBus := eventbus.InitEventBus()
	queryBus := querybus.InitQueryBus()
	deduplicator := messaging.NewDefaultRedisDeduplicator(redisClient)
	mutexService := distributedsync.NewRedisMutexService(redisClient, appLogger)
	uuidProvider := utils.NewRandomUUIDProvider()

	return &CommonServices{
		Config:       cfg,
		Logger:       appLogger,
		RedisClient:  redisClient,
		EventBus:     eventBus,
		QueryBus:     queryBus,
		Deduplicator: deduplicator,
		Mutex:        mutexService,
		Router:       &router,
		UUIDProvider: uuidProvider,
		TimeProvider: timeProvider,
	}
}

func MustInitCommonServicesWithEnvFiles(ctx context.Context, envFiles ...string) *CommonServices {
	err := godotenv.Overload(envFiles...)
	if err != nil {
		panic(err)
	}

	return MustInitCommonServices(ctx)
}
