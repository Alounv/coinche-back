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

func TestGameCreation(test *testing.T) {
	assert := assert.New(test)
	dbName := "testdb"
	env.LoadEnv("../.env")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")

	db := testUtils.CreateDb(connectionInfo, dbName)

	MockGameService := NewGameServiceWithData(db)

	test.Run("get a game", func(test *testing.T) {
		want := "GAME ONE"

		got := MockGameService.GetGame(1)

		assert.Equal(want, got.Name)
		assert.Equal(1, got.Id)
	})

	test.Run("list all games", func(test *testing.T) {
		want := []app.Game{
			{Name: "GAME ONE", Id: 1},
			{Name: "GAME TWO", Id: 2},
		}

		got := MockGameService.ListGames()

		assert.Equal(want, got)
	})

	test.Cleanup(func() {
		testUtils.DropDb(connectionInfo, dbName, db)
	})
}

func NewGameServiceWithData(db *sqlx.DB) *GameService {
	store := NewGameServiceFromDb(db)

	store.CreatePlayerTableIfNeeded()
	store.CreateGames([]app.Game{
		{Name: "GAME ONE", Id: 1},
		{Name: "GAME TWO", Id: 2},
	})

	return store
}
