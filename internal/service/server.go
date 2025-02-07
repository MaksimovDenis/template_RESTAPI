package service

import (
	"github.com/gin-gonic/gin"
)

type ServerService struct{}

func NewServerService() *ServerService {
	return &ServerService{}
}

func (serv *ServerService) CheckService(ctx *gin.Context) error {
	return nil
}
