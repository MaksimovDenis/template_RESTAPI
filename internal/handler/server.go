package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (hdl *Handler) CheckServer(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}
