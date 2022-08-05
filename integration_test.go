package main

import (
	"coinche/api"
	"coinche/domain"
	repository "coinche/repository"
	"coinche/usecases"
	"coinche/utilities"
	testUtilities "coinche/utilities/test"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	db             *sqlx.DB
	connectionInfo string
	dbName         string
	router         *gin.Engine
	gameUsecases   *usecases.GameUsecases
	server1        *httptest.Server
	server2        *httptest.Server
	server3        *httptest.Server
	server4        *httptest.Server
	connection1    *websocket.Conn
	connection2    *websocket.Conn
	connection3    *websocket.Conn
	connection4    *websocket.Conn
	hub            *api.Hub
}

func (s *IntegrationTestSuite) SetupSuite() {
	utilities.LoadEnv("")
	s.connectionInfo = os.Getenv("SQLX_POSTGRES_INFO")
	s.dbName = "testdb"

	s.db = testUtilities.CreateDb(s.connectionInfo, s.dbName)

	gameRepository, err := repository.NewGameRepositoryFromDb(s.db)
	if err != nil {
		s.T().Fatal(err)
	}

	s.gameUsecases = &usecases.GameUsecases{Repo: gameRepository}

	s.router, s.hub = api.SetupRouter(s.gameUsecases)
}

func (s *IntegrationTestSuite) TestCreateGame() {
	test := s.T()
	assert := assert.New(test)
	response := httptest.NewRecorder()

	test.Run("create game", func(test *testing.T) {
		s.router.ServeHTTP(httptest.NewRecorder(), testUtilities.NewCreateGameRequest(test, "NEW GAME"))

		assert.Equal(http.StatusOK, response.Code)
	})

	test.Run("get game", func(test *testing.T) {
		s.router.ServeHTTP(response, testUtilities.NewGetGameRequest(test, 1))

		got := testUtilities.DecodeToGame(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal("NEW GAME", got.Name)
		assert.Equal(1, got.ID)
		assert.Equal(map[string]domain.Player{}, got.Players)
		assert.IsType(time.Time{}, got.CreatedAt)
	})

	test.Run("list game", func(test *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/games/all", nil)
		s.router.ServeHTTP(response, request)

		got := testUtilities.DecodeToGames(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(1, len(got))
		assert.Equal("NEW GAME", got[0].Name)
		assert.Equal(1, got[0].ID)
		assert.IsType(time.Time{}, got[0].CreatedAt)
	})

	test.Run("join game", func(test *testing.T) {
		s.server1, s.connection1 = api.NewGameWebSocketServer(test, 1, "P1", s.hub)
		receivedGame := api.ReceiveGameOrFatal(s.connection1, test)

		assert.IsType(domain.Game{}, receivedGame)

		assert.Equal(map[string]domain.Player{"P1": {}}, receivedGame.Players)
	})

	test.Run("leave unstarted game", func(test *testing.T) {
		err := api.SendMessage(s.connection1, "leave")
		if err != nil {
			test.Fatal(err)
		}

		message := api.ReceiveMessageOrFatal(s.connection1, test)

		assert.Equal("Has left the game", message)

		s.router.ServeHTTP(response, testUtilities.NewGetGameRequest(test, 1))
		got := testUtilities.DecodeToGame(response.Body, test)

		assert.Equal(map[string]domain.Player{}, got.Players)
	})

	test.Run("join back unstarted game", func(test *testing.T) {
		s.server1, s.connection1 = api.NewGameWebSocketServer(test, 1, "P1", s.hub)
		got := api.ReceiveGameOrFatal(s.connection1, test)

		assert.IsType(domain.Game{}, got)

		assert.Equal(map[string]domain.Player{"P1": {}}, got.Players)
	})

	test.Run("other players join", func(test *testing.T) {
		s.server2, s.connection2 = api.NewGameWebSocketServer(test, 1, "P2", s.hub)
		s.server3, s.connection3 = api.NewGameWebSocketServer(test, 1, "P3", s.hub)
		s.server4, s.connection4 = api.NewGameWebSocketServer(test, 1, "P4", s.hub)
		_ = api.ReceiveGameOrFatal(s.connection1, test)
		_ = api.ReceiveGameOrFatal(s.connection1, test)
		_ = api.ReceiveGameOrFatal(s.connection1, test)
		_ = api.ReceiveGameOrFatal(s.connection2, test)
		_ = api.ReceiveGameOrFatal(s.connection2, test)
		_ = api.ReceiveGameOrFatal(s.connection2, test)
		_ = api.ReceiveGameOrFatal(s.connection3, test)
		_ = api.ReceiveGameOrFatal(s.connection3, test)
		got := api.ReceiveGameOrFatal(s.connection4, test)

		assert.IsType(domain.Game{}, got)

		assert.Equal(map[string]domain.Player{"P1": {}, "P2": {}, "P3": {}, "P4": {}}, got.Players)
	})

	test.Run("start game should return error if team not ready", func(test *testing.T) {
		err := api.SendMessage(s.connection1, "start")
		if err != nil {
			test.Fatal(err)
		}

		got := api.ReceiveMessageOrFatal(s.connection1, test)

		assert.Equal("Could not start the game: TEAMS ARE NOT EQUAL", got)
	})

	test.Run("join team", func(test *testing.T) {
		err := api.SendMessage(s.connection1, "joinTeam: Odd")
		if err != nil {
			test.Fatal(err)
		}
		err = api.SendMessage(s.connection2, "joinTeam: Even")
		if err != nil {
			test.Fatal(err)
		}
		err = api.SendMessage(s.connection3, "joinTeam: Odd")
		if err != nil {
			test.Fatal(err)
		}
		err = api.SendMessage(s.connection4, "joinTeam: Even")
		if err != nil {
			test.Fatal(err)
		}

		_ = api.ReceiveGameOrFatal(s.connection1, test)
		_ = api.ReceiveGameOrFatal(s.connection1, test)
		_ = api.ReceiveGameOrFatal(s.connection1, test)
		got := api.ReceiveGameOrFatal(s.connection1, test)

		assert.Equal(map[string]domain.Player{"P1": {Team: "Odd"}, "P2": {Team: "Even"}, "P3": {Team: "Odd"}, "P4": {Team: "Even"}}, got.Players)
	})
}

func (s *IntegrationTestSuite) TearDownSuite() {
	testUtilities.DropDb(s.connectionInfo, s.dbName, s.db)
	s.server1.Close()
	s.server2.Close()
	s.server3.Close()
	s.server4.Close()
	s.connection1.Close()
	s.connection2.Close()
	s.connection3.Close()
	s.connection4.Close()
}

func TestIntegrationSuite(test *testing.T) {
	suite.Run(test, new(IntegrationTestSuite))
}
