package handler

import (
	"net/http"
	"templates_new/pkg/protocol/oapi"

	"github.com/gin-gonic/gin"
)

func (hdl *Handler) SignIn(ctx *gin.Context) {
	var user oapi.SignInJSONBody

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse job create object."})

		return
	}

	res, err := hdl.serverService.Autorization.SignIn(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job."})

		return
	}

	ctx.JSON(http.StatusCreated, res)
}
func (hdl *Handler) LogIn(ctx *gin.Context) {}
