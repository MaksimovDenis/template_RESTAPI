package handler

import (
	"templates_new/internal/service"
	"templates_new/pkg/protocol/oapi"
	"templates_new/pkg/token"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	FileUploadBufferSize       = 512e+6 // 512MB for now
	ServerShutdownDefaultDelay = 5 * time.Second
)

type Handler struct {
	appService service.Service
	tokenMaker *token.JWTMaker
	log        zerolog.Logger
}

func NewHandler(appService service.Service, tokenMaker token.JWTMaker, log zerolog.Logger) *Handler {
	return &Handler{
		appService: appService,
		tokenMaker: &tokenMaker,
		log:        log,
	}
}

func (hdl *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.MaxMultipartMemory = FileUploadBufferSize

	tokenMaker := hdl.tokenMaker

	oapi.RegisterHandlersWithOptions(router, hdl, oapi.GinServerOptions{
		BaseURL: "/api",
		Middlewares: []oapi.MiddlewareFunc{
			GetAuthMiddlewareFunc(tokenMaker),
		},
	})

	return router
}
