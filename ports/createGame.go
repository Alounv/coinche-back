package ports

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) CreateGame(context *gin.Context) {
	name := context.Query("name")

	id := server.Store.CreateGame(name)
	context.JSON(http.StatusAccepted, id)
}
