package testUtilities

import (
	"bytes"
	"coinche/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"testing"

	epg "github.com/fergusstrange/embedded-postgres"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

const connectionInfo = "user=postgres password=postgres port=5432"

func NewCreateGameRequest(test *testing.T, name string) *http.Request {
	route := fmt.Sprintf("/games/create?name=%s", url.QueryEscape(name))
	return GetNewRequest(test, route, http.MethodPost)
}

func NewGetGameRequest(test *testing.T, gameID int) *http.Request {
	route := fmt.Sprintf("/games/%d", gameID)
	return GetNewRequest(test, route, http.MethodGet)
}

func NewJoinGameRequest(test *testing.T, gameID int, playerName string) *http.Request {
	route := fmt.Sprintf("/games/%d/join?playerName=%s", gameID, url.QueryEscape(playerName))
	return GetNewRequest(test, route, http.MethodGet)
}

func GetNewRequest(test *testing.T, route string, method string) *http.Request {
	request, err := http.NewRequest(method, route, nil)
	if err != nil {
		test.Fatal(err)
	}
	return request
}

func CreateDb(dbName string) (*sqlx.DB, *epg.EmbeddedPostgres) {
	postgres := epg.NewDatabase()

	err := postgres.Start()
	if err != nil {
		fmt.Println("Could not start postgres: ", err)
	}

	userDb := sqlx.MustOpen("pgx", connectionInfo)

	_, err = userDb.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		fmt.Println("Could not create DB, tries to  drop before creation", err)
		userDb.MustExec("DROP DATABASE " + dbName)
		userDb.MustExec("CREATE DATABASE " + dbName)
	}
	userDb.Close()

	db := sqlx.MustOpen("pgx", connectionInfo+" dbname="+dbName)
	return db, postgres
}

func DropDb(postgres *epg.EmbeddedPostgres, dbName string, db *sqlx.DB) {
	db.Close()

	userDb := sqlx.MustOpen("pgx", connectionInfo)
	userDb.MustExec("DROP DATABASE " + dbName)
	userDb.Close()

	err := postgres.Stop()
	fmt.Println(err)
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

func GetSortedPlayersNameByOrder(players map[string]domain.Player) []string {
	type nameOrder struct {
		name  string
		order int
	}

	ss := []nameOrder{}
	for name, player := range players {
		ss = append(ss, nameOrder{name, player.Order})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].order < ss[j].order
	})

	sortedPlayerNames := []string{}
	for _, pair := range ss {
		sortedPlayerNames = append(sortedPlayerNames, pair.name)
	}
	return sortedPlayerNames
}

func GetPlayerIndexFromNameOrFatal(name string, test *testing.T) int {
	playerNb, err := strconv.Atoi(name[1:2])
	if err != nil {
		test.Fatal(err)
	}
	return playerNb - 1
}
