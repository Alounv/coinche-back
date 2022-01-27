package adapters

import (
	"coinche/app"
	"coinche/utilities/env"
	testUtils "coinche/utilities/test"
	"os"
	"testing"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGameGetListUpdate(test *testing.T) {
	assert := assert.New(test)
	dbName := "testdb"
	env.LoadEnv("../.env")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")

	db := testUtils.CreateDb(connectionInfo, dbName)

	MockGameService := NewGameServiceWithData(db)

	test.Run("get a game", func(test *testing.T) {
		want := app.Game{Name: "GAME ONE", Id: 1, Players: []string{}}

		got := MockGameService.GetGame(1)

		assert.Equal(want, got)
	})

	test.Run("list all games", func(test *testing.T) {
		want := []app.Game{
			{Name: "GAME ONE", Id: 1, Players: []string{}},
			{Name: "GAME TWO", Id: 2, Players: []string{"P1", "P2"}},
		}

		got := MockGameService.ListGames()

		assert.Equal(want, got)
	})

	test.Run("update a game", func(test *testing.T) {
		want := []string{"P1", "P2"}

		MockGameService.UpdateGame(2, want)
		got := MockGameService.GetGame(2).Players

		assert.Equal(want, got)
	})

	test.Cleanup(func() {
		testUtils.DropDb(connectionInfo, dbName, db)
	})
}

func NewGameServiceWithData(db *sqlx.DB) *dbGameService {
	dbGameService := NewDbGameServiceFromDb(db)

	dbGameService.CreateGames([]app.Game{
		{Name: "GAME ONE", Id: 1},
		{Name: "GAME TWO", Id: 2, Players: []string{"P1", "P2"}},
	})

	return dbGameService
}
