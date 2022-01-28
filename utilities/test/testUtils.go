package testUtils

import (
	"bytes"
	"coinche/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"
)

func NewCreateGameRequest(name string) *http.Request {
	route := fmt.Sprintf("/games/create?name=%s", name)
	return GetNewRequest(route, http.MethodPost)
}

func NewGetGameRequest(id int) *http.Request {
	route := fmt.Sprintf("/games/%d", id)
	return GetNewRequest(route, http.MethodGet)
}

func GetNewRequest(route string, method string) *http.Request {
	request, err := http.NewRequest(method, route, nil)
	if err != nil {
		panic(err)
	}
	return request
}

func CreateDb(connectionInfo string, dbName string) *sqlx.DB {
	userDb := sqlx.MustOpen("pgx", connectionInfo)
	_, err := userDb.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		fmt.Print(err)
		userDb.MustExec("DROP DATABASE " + dbName)
		userDb.MustExec("CREATE DATABASE " + dbName)
	}
	userDb.Close()

	db := sqlx.MustOpen("pgx", connectionInfo+" dbname="+dbName)
	return db
}

func DropDb(connectionInfo string, dbName string, db *sqlx.DB) {
	db.Close()

	userDb := sqlx.MustOpen("pgx", connectionInfo)
	userDb.MustExec("DROP DATABASE " + dbName)
	userDb.Close()
}

func DecodeToGames(buf *bytes.Buffer, test *testing.T) []domain.Game {
	var got []domain.Game
	err := json.NewDecoder(buf).Decode(&got)
	if err != nil {
		test.Fatalf("Unable to parse response from gameAPIs %q into %q, '%v'", buf, "slice of Game", err)
	}
	return got
}

func DecodeToGame(buf *bytes.Buffer, test *testing.T) domain.Game {
	var got domain.Game
	err := json.NewDecoder(buf).Decode(&got)
	if err != nil {
		test.Fatalf("Unable to parse response from gameAPIs %q into %q, '%v'", buf, "Game", err)
	}
	return got
}
