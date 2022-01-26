package ports

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) CreateGame(context *gin.Context) {
	name := context.Param("name")

	server.Store.CreateGame(name)
	context.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}
