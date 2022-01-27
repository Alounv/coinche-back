package main

import (
	"coinche/adapters"
	"coinche/ports"
	"coinche/utilities/env"
	testUtils "coinche/utilities/test"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreatingAndGettingGames(test *testing.T) {
	assert := assert.New(test)

	env.LoadEnv("")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")
	dbName := "testdb"

	db := testUtils.CreateDb(connectionInfo, dbName)
	mockGameService := adapters.NewGameServiceFromDb(db)

	router := ports.SetupRouter(mockGameService)

	router.ServeHTTP(httptest.NewRecorder(), testUtils.NewCreateGameRequest("NEW GAME"))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, testUtils.NewGetGameRequest(1))
	got := testUtils.DecodeToGame(response.Body, test)

	assert.Equal(http.StatusOK, response.Code)
	assert.Equal("NEW GAME", got.Name)
	assert.Equal(1, got.Id)
	assert.IsType(time.Time{}, got.CreatedAt)

	test.Cleanup(func() {
		testUtils.DropDb(connectionInfo, dbName, db)
	})
}
