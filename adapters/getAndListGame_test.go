package adapters

import (
	"coinche/app"
	"coinche/utilities/env"
	testUtils "coinche/utilities/test"
	"os"
	"testing"
	"time"

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
			{Name: "GAME ONE", Id: 1, CreatedAt: time.Date(2009, 1, 1, 12, 0, 0, 0, time.UTC)},
			{Name: "GAME TWO", Id: 2, CreatedAt: time.Date(2009, 1, 2, 12, 0, 0, 0, time.UTC)},
		}

		got := MockGameService.ListGames()

		assert.Equal(want, got)
	})

	test.Cleanup(func() {
		testUtils.DropDb(connectionInfo, dbName, db)
	})
}

func NewGameServiceWithData(db *sqlx.DB) *dbGameService {
	store := NewDbGameServiceFromDb(db)

	store.CreateGames([]app.Game{
		{Name: "GAME ONE", Id: 1, CreatedAt: time.Date(2009, 1, 1, 12, 0, 0, 0, time.UTC)},
		{Name: "GAME TWO", Id: 2, CreatedAt: time.Date(2009, 1, 2, 12, 0, 0, 0, time.UTC)},
	})

	return store
}
