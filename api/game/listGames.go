package gameApi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (gameAPIs *GameAPIs) ListGames(context *gin.Context) {
	games := gameAPIs.Store.ListGames()

	context.JSON(http.StatusOK, games)
}
