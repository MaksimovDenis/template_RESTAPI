package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (hdl *Handler) probe(ctx *gin.Context) {
	err := hdl.serverService.Probe(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}

func (hdl *Handler) CheckServer(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}
