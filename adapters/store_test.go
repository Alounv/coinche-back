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

func TestDB(test *testing.T) {
	assert := assert.New(test)
	dbName := "testdb"
	env.LoadEnv("../.env")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")

	db := testUtils.CreateDb(connectionInfo, dbName)

	MockGameService := NewGameServiceWithData(db)

	test.Run("get a game", func(test *testing.T) {
		want := "GAME ONE"

		got := MockGameService.GetAGame(1)

		assert.Equal(want, got.Name)
		assert.NotNil(got.Id)
	})

	test.Run("create a game", func(test *testing.T) {
		newName := "NEW GAME"

		newId := MockGameService.CreateAGame(newName)
		got := MockGameService.GetAGame(newId)

		assert.Equal(newName, got.Name)
		assert.NotNil(newId, got.Id)
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
