package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/tclaudel/hexagonal-architecture-golang/internal/adapters/primary/user/server"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/adapters/secondary/user/mongo"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/config"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/services"
)

const shutdownTimeout = 5 * time.Second

func main() {
	ctx := context.Background()

	cfg := config.NewConfig()

	serverContainer := newServiceContainer(cfg)

	start(ctx, serverContainer)
	waitForInterruptSignal()
	shutdown(ctx, serverContainer)
}

func start(ctx context.Context, sc *serviceContainer) {
	go func() {
		err := sc.getUserServer(ctx).ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
}

func waitForInterruptSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func shutdown(ctx context.Context, sc *serviceContainer) {
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	err := sc.getUserServer(ctx).Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}

type serviceContainer struct {
	cfg                 *config.Config
	userMongoRepository *mongo.UserRepository
	userUseCase         *services.UserUseCase
	userServer          *server.Server
}

func newServiceContainer(cfg *config.Config) *serviceContainer {
	return &serviceContainer{
		cfg: cfg,
	}
}

func (sc serviceContainer) getUserMongoRepository(ctx context.Context) *mongo.UserRepository {
	if sc.userMongoRepository == nil {
		mongoUserRepository, err := mongo.NewMongoRepository(ctx, sc.cfg.MongoURI(), sc.cfg.MongoDBName())
		if err != nil {
			panic(err)
		}

		sc.userMongoRepository = mongoUserRepository
	}

	return sc.userMongoRepository
}

func (sc serviceContainer) getUserUseCase(ctx context.Context) *services.UserUseCase {
	if sc.userUseCase == nil {
		sc.userUseCase = services.NewUserUseCase(sc.getUserMongoRepository(ctx))
	}

	return sc.userUseCase
}

func (sc serviceContainer) getUserServer(ctx context.Context) *server.Server {
	if sc.userServer == nil {
		userServer, err := server.NewServer(ctx, server.Params{
			UserUseCase: sc.getUserUseCase(ctx),
			Port:        sc.cfg.UserServerPort(),
		})
		if err != nil {
			panic(err)
		}

		sc.userServer = userServer
	}

	return sc.userServer
}
