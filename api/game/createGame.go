package gameapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (gameAPIs *GameAPIs) CreateGame(context *gin.Context) {
	name := context.Query("name")
	creatorName := context.Query("creator")

	id := gameAPIs.GameService.CreateGame(name, creatorName)
	context.JSON(http.StatusAccepted, id)
}
