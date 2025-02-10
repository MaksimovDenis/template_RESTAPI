package service

import (
	"github.com/gin-gonic/gin"
)

type Server interface {
	CheckService(ctx *gin.Context) error
}

type ServerService struct{}

func newServerService() *ServerService {
	return &ServerService{}
}

func (serv *ServerService) CheckService(ctx *gin.Context) error {
	return nil
}
