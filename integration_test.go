package main

import (
	"coinche/api"
	"coinche/app"
	gameRepo "coinche/repository/game"
	"coinche/utilities/env"
	testUtils "coinche/utilities/test"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func SetupTest() (*gin.Engine, *sqlx.DB, string, string) {
	env.LoadEnv("")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")
	dbName := "testdb"

	db := testUtils.CreateDb(connectionInfo, dbName)

	gameRepo := gameRepo.NewGameRepoFromDb(db)
	gameService := &app.GameService{GameRepo: gameRepo}

	router := api.SetupRouter(gameService)
	return router, db, connectionInfo, dbName
}

func TestCreateAndGetGames(test *testing.T) {
	assert := assert.New(test)
	router, db, connectionInfo, dbName := SetupTest()

	response := httptest.NewRecorder()
	router.ServeHTTP(httptest.NewRecorder(), testUtils.NewCreateGameRequest("NEW GAME"))

	router.ServeHTTP(response, testUtils.NewGetGameRequest(1))

	router.ServeHTTP(response, testUtils.NewJoinGameRequest(1, "player1"))

	got := testUtils.DecodeToGame(response.Body, test)

	assert.Equal(http.StatusOK, response.Code)
	assert.Equal("NEW GAME", got.Name)
	assert.Equal(1, got.Id)
	assert.IsType(time.Time{}, got.CreatedAt)

	test.Cleanup(func() {
		testUtils.DropDb(connectionInfo, dbName, db)
	})
}
