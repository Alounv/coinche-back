package gameapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateGame is exported for other packages such as repository and api
func (gameAPIs *GameAPIs) CreateGame(context *gin.Context) {
	name := context.Query("name")

	id := gameAPIs.Store.CreateGame(name)
	context.JSON(http.StatusAccepted, id)
}
