package handler

import (
	"templates_new/internal/service"
	"templates_new/pkg/protocol/oapi"
	"templates_new/pkg/token"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	FileUploadBufferSize       = 512e+6 // 512MB for now
	ServerShutdownDefaultDelay = 5 * time.Second
)

type Handler struct {
	serverService service.Service
	tokenMaker    *token.JWTMaker
}

func NewHandler(appService service.Service, secretKey string) *Handler {
	return &Handler{
		serverService: appService,
		tokenMaker:    token.NewJWTMaker(secretKey),
	}
}

func (hdl *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.MaxMultipartMemory = FileUploadBufferSize

	tokenMaker := hdl.tokenMaker

	adminRoutes := router.Group("/api/admin")
	adminRoutes.Use(GetAdminMiddlewareFunc(tokenMaker))

	userRoutes := router.Group("/api/user")
	userRoutes.Use(GetAuthMiddlewareFunc(tokenMaker))

	oapi.RegisterHandlersWithOptions(router, hdl, oapi.GinServerOptions{
		BaseURL: "/api",
	})

	return router
}
