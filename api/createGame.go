package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (gameAPIs *GameAPIs) CreateGame(context *gin.Context) {
	name := context.Query("name")

	gameID, err := gameAPIs.Usecases.CreateGame(name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusAccepted, gameID)
}
