package gameapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListGames is exported for other packages such as repository and api
func (gameAPIs *GameAPIs) ListGames(context *gin.Context) {
	games := gameAPIs.Store.ListGames()

	context.JSON(http.StatusOK, games)
}
