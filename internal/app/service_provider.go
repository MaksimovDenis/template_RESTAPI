package app

import (
	"context"
	"os"
	db "templates_new/internal/client"
	"templates_new/internal/client/db/pg"
	"templates_new/internal/client/db/transaction"
	"templates_new/internal/config"
	"templates_new/internal/handler"
	"templates_new/internal/repository"
	"templates_new/internal/service"
	"templates_new/pkg/token"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type serviceProvider struct {
	pgConfig     config.PGConfig
	serverConfig config.ServerConfig
	tokenConfig  config.TokenConfig

	dbClient      db.Client
	txManager     db.TxManager
	appRepository *repository.Repository

	appService *service.Service

	tokenMaker *token.JWTMaker

	log zerolog.Logger

	handler *handler.Handler
}

func newServiceProvider() *serviceProvider {
	srv := &serviceProvider{}
	srv.log = srv.initLogger()

	return srv
}

func (srv *serviceProvider) PGConfig() config.PGConfig {
	if srv.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get pg config")
		}

		srv.pgConfig = cfg
	}

	return srv.pgConfig
}

func (srv *serviceProvider) ServerConfig() config.ServerConfig {
	if srv.serverConfig == nil {
		cfg, err := config.NewServerConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get server config")
		}

		srv.serverConfig = cfg
	}

	return srv.serverConfig
}

func (srv *serviceProvider) TokenConfig() config.TokenConfig {
	if srv.tokenConfig == nil {
		cfg, err := config.NewSecretKey()
		if err != nil {
			log.Fatal().Err(err).Msg("failed dto get secret key config")
		}

		srv.tokenConfig = cfg
	}

	return srv.tokenConfig
}

func (srv *serviceProvider) DBClient(ctx context.Context) db.Client {
	if srv.dbClient == nil {
		cl, err := pg.New(ctx, srv.PGConfig().DSN())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create db client")
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("ping error")
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
		srv.appRepository = repository.NewRepository(
			srv.DBClient(ctx),
			srv.log.With().Str("module", "repository").Logger(),
		)
	}

	return srv.appRepository
}

func (srv *serviceProvider) AppService(ctx context.Context) *service.Service {
	if srv.appService == nil {
		srv.appService = service.NewService(
			*srv.AppRepository(ctx),
			srv.TxManager(ctx),
			*srv.TokenMaker(ctx),
			srv.log.With().Str("module", "service").Logger(),
		)
	}

	return srv.appService
}

func (srv *serviceProvider) TokenMaker(ctx context.Context) *token.JWTMaker {
	if srv.tokenMaker == nil {
		srv.tokenMaker = token.NewJWTMaker(
			srv.TokenConfig().SecretKey(),
		)
	}

	return srv.tokenMaker
}

func (srv *serviceProvider) initLogger() zerolog.Logger {
	logFile, err := os.OpenFile("./internal/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open log file")
	}

	logLevel, err := zerolog.ParseLevel("debug")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	}

	// Записываем логи и в файл, и в консоль
	multiWriter := zerolog.MultiLevelWriter(os.Stdout, logFile)

	logger := zerolog.New(multiWriter).Level(logLevel).With().Timestamp().Logger()
	return logger
}

func (srv *serviceProvider) AppHandler(ctx context.Context) *handler.Handler {
	if srv.handler == nil {
		srv.handler = handler.NewHandler(
			*srv.AppService(ctx),
			*srv.TokenMaker(ctx),
			srv.log.With().Str("module", "api").Logger(),
		)
	}
	return srv.handler
}
