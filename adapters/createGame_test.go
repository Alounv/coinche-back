package adapters

import (
	"coinche/utilities/env"
	testUtils "coinche/utilities/test"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestCreateGame(test *testing.T) {
	assert := assert.New(test)
	dbName := "testdb"
	env.LoadEnv("../.env")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")

	db := testUtils.CreateDb(connectionInfo, dbName)

	MockGameService := NewGameServiceFromDb(db)

	test.Run("create a game", func(test *testing.T) {
		newName := "NEW GAME"

		newId := MockGameService.CreateGame(newName)
		got := MockGameService.GetGame(newId)

		assert.Equal(newName, got.Name)
		assert.Equal(newId, got.Id)
		assert.IsType(time.Time{}, got.CreatedAt)
	})

	test.Cleanup(func() {
		testUtils.DropDb(connectionInfo, dbName, db)
	})
}
