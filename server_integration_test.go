package main

import (
	"coinche/adapters"
	"coinche/app"
	"coinche/ports"
	"coinche/utilities/env"
	testUtils "coinche/utilities/test"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSettingAndGettingScores(test *testing.T) {
	assert := assert.New(test)

	env.LoadEnv("")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")
	dbName := "testdb"

	db := testUtils.CreateDb(connectionInfo, dbName)
	mockGameService := NewGameServiceWithData(db)

	router := ports.SetupRouter(mockGameService)

	router.ServeHTTP(httptest.NewRecorder(), testUtils.NewPOSTGameRequest("NEW GAME"))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, testUtils.NewGETGameRequest(1))


	assert.Equal(http.StatusOK, response.Code)
	assert.Equal("{\"Name\":\"GAME ONE\",\"Id\":1}", response.Body.String())

	test.Cleanup(func() {
		testUtils.DropDb(connectionInfo, dbName, db)
	})
}

func NewGameServiceWithData (db *sqlx.DB) *adapters.GameService {
	store := adapters.NewGameServiceFromDb(db)

	store.CreatePlayerTableIfNeeded()
	store.CreateGames([]app.Game{
		{Name: "GAME ONE", Id: 1},
		{Name: "GAME TWO", Id: 2},
	})

	return store
}