package testUtils

import (
	"bytes"
	"coinche/app"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"
)

func NewPOSTGameRequest(name string) *http.Request {
	route := fmt.Sprintf("/games/%s", name)
	return getNewRequest(route, http.MethodPost)
}

func NewGETGameRequest(id int) *http.Request {
	route := fmt.Sprintf("/games/%d", id)
	return getNewRequest(route, http.MethodGet)
}

func getNewRequest(route string, method string) *http.Request {
	request, err := http.NewRequest(method, route, nil)
	if err != nil {
		panic(err)
	}
	return request
}

func CreateDb (connectionInfo string, dbName string) *sqlx.DB {
	userDb := sqlx.MustOpen("pgx", connectionInfo) 
	_, err := userDb.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		fmt.Print("Hello")
		fmt.Print(err)
	}
	userDb.Close()

	db := sqlx.MustOpen("pgx", connectionInfo + " dbname=" + dbName)
	return db
}

func DropDb (connectionInfo string, dbName string, db *sqlx.DB) {
	db.Close()

	userDb := sqlx.MustOpen("pgx", connectionInfo) 
	userDb.MustExec("DROP DATABASE " + dbName)
	userDb.Close()
}

func DecodeToGames (buf *bytes.Buffer, test *testing.T) []app.Game {
	var got []app.Game
	err := json.NewDecoder(buf).Decode(&got)
	if err != nil {
		test.Fatalf("Unable to parse response from server %q into %q, '%v'", buf, "slice of Game", err)
	}
	return got
}

func DecodeToGame (buf *bytes.Buffer, test *testing.T) app.Game {
	var got app.Game
	err := json.NewDecoder(buf).Decode(&got)
	if err != nil {
		test.Fatalf("Unable to parse response from server %q into %q, '%v'", buf, "Game", err)
	}
	return got
}