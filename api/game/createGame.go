package gameapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (gameAPIs *GameAPIs) CreateGame(context *gin.Context) {
	name := context.Query("name")

	id := gameAPIs.GameService.CreateGame(name)
	context.JSON(http.StatusAccepted, id)
}
