package ports

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) ListGames(context *gin.Context) {
	games := server.Store.ListGames()

	context.JSON(http.StatusOK, games)
}
