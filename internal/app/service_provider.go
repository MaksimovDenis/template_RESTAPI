package app

import (
	"context"
	"log"
	db "templates_new/internal/client"
	"templates_new/internal/client/db/pg"
	"templates_new/internal/client/db/transaction"
	"templates_new/internal/config"
	"templates_new/internal/handler"
	"templates_new/internal/repository"
	"templates_new/internal/service"
)

type serviceProvider struct {
	pgConfig     config.PGConfig
	serverConfig config.ServerConfig
	tokenConfig  config.TokenConfig

	dbClient      db.Client
	txManager     db.TxManager
	appRepository *repository.Repository

	appService *service.Service

	handler *handler.Handler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (srv *serviceProvider) PGConfig() config.PGConfig {
	if srv.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		srv.pgConfig = cfg
	}

	return srv.pgConfig
}

func (srv *serviceProvider) ServerConfig() config.ServerConfig {
	if srv.serverConfig == nil {
		cfg, err := config.NewServerConfig()
		if err != nil {
			log.Fatalf("failed to get server config: %s", err.Error())
		}

		srv.serverConfig = cfg
	}

	return srv.serverConfig
}

func (srv *serviceProvider) TokenConfig() config.TokenConfig {
	if srv.tokenConfig == nil {
		cfg, err := config.NewSecretKey()
		if err != nil {
			log.Fatalf("failed dto get secret key config: %s", err.Error())
		}

		srv.tokenConfig = cfg
	}

	return srv.tokenConfig
}

func (srv *serviceProvider) DBClient(ctx context.Context) db.Client {
	if srv.dbClient == nil {
		cl, err := pg.New(ctx, srv.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}

		srv.dbClient = cl
	}

	return srv.dbClient
}

func (srv *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if srv.txManager == nil {
		srv.txManager = transaction.NewTransactionsManager(srv.DBClient(ctx).DB())
	}

	return srv.txManager
}

func (srv *serviceProvider) AppRepository(ctx context.Context) *repository.Repository {
	if srv.appRepository == nil {
		srv.appRepository = repository.NewRepository(srv.DBClient(ctx))
	}

	return srv.appRepository
}

func (srv *serviceProvider) ServerService(ctx context.Context) service.Service {
	if srv.appService == nil {
		srv.appService = service.NewService(
			*srv.AppRepository(ctx),
			srv.TxManager(ctx),
		)
	}

	return *srv.appService
}

func (srv *serviceProvider) ServiceHandler(ctx context.Context) *handler.Handler {
	if srv.handler == nil {
		srv.handler = handler.NewHandler(
			srv.ServerService(ctx),
			srv.TokenConfig().SecretKey(),
		)
	}

	return srv.handler
}
