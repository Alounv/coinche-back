package main

import (
	"coinche/api"
	"coinche/domain"
	gamerepo "coinche/repository/game"
	"coinche/usecases"
	"coinche/utilities"
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
	server         *httptest.Server
	connection     *websocket.Conn
	hub            *api.Hub
}

func (s *IntegrationTestSuite) SetupSuite() {
	utilities.LoadEnv("")
	s.connectionInfo = os.Getenv("SQLX_POSTGRES_INFO")
	s.dbName = "testdb"

	s.db = utilities.CreateDb(s.connectionInfo, s.dbName)

	gameRepository, err := gamerepo.NewGameRepositoryFromDb(s.db)
	if err != nil {
		s.T().Fatal(err)
	}

	s.gameUsecases = &usecases.GameUsecases{Repo: gameRepository}

	s.router, s.hub = api.SetupRouter(s.gameUsecases)
}

func (s *IntegrationTestSuite) TestCreateGame() {
	test := s.T()
	assert := assert.New(test)

	s.router.ServeHTTP(httptest.NewRecorder(), utilities.NewCreateGameRequest(test, "NEW GAME"))

	response := httptest.NewRecorder()
	assert.Equal(http.StatusOK, response.Code)
}

func (s *IntegrationTestSuite) TestGetGame() {
	test := s.T()
	assert := assert.New(test)

	response := httptest.NewRecorder()
	s.router.ServeHTTP(response, utilities.NewGetGameRequest(test, 1))

	got := utilities.DecodeToGame(response.Body, test)

	assert.Equal(http.StatusOK, response.Code)
	assert.Equal("NEW GAME", got.Name)
	assert.Equal(1, got.ID)
	assert.Equal(map[string]domain.Player{}, got.Players)
	assert.IsType(time.Time{}, got.CreatedAt)
}

func (s *IntegrationTestSuite) TestListGames() {
	test := s.T()
	assert := assert.New(test)
	request, _ := http.NewRequest(http.MethodGet, "/games/all", nil)
	response := httptest.NewRecorder()

	s.router.ServeHTTP(response, request)

	got := utilities.DecodeToGames(response.Body, test)

	assert.Equal(http.StatusOK, response.Code)
	assert.Equal(1, len(got))
	assert.Equal("NEW GAME", got[0].Name)
	assert.Equal(1, got[0].ID)
	assert.IsType(time.Time{}, got[0].CreatedAt)
}

func (s *IntegrationTestSuite) TestJoinGame() {
	test := s.T()
	assert := assert.New(test)
	response := httptest.NewRecorder()

	s.server, s.connection = api.NewGameWebSocketServer(test, s.gameUsecases, 1, "player", s.hub)
	receivedGame, _ := api.ReceiveGame(s.connection)

	assert.IsType(domain.Game{}, receivedGame)

	s.router.ServeHTTP(response, utilities.NewGetGameRequest(test, 1))
	got := utilities.DecodeToGame(response.Body, test)

	assert.Equal(map[string]domain.Player{"player": {}}, got.Players)
}

func (s *IntegrationTestSuite) TestLeaveUnstartedGame() {
	test := s.T()
	assert := assert.New(test)
	response := httptest.NewRecorder()

	err := api.SendMessage(s.connection, "leave")
	utilities.FatalIfErr(err, test)

	message, _ := api.ReceiveMessage(s.connection)

	assert.Equal("Has left the game", message)

	s.router.ServeHTTP(response, utilities.NewGetGameRequest(test, 1))
	got := utilities.DecodeToGame(response.Body, test)

	assert.Equal(map[string]domain.Player{}, got.Players)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	utilities.DropDb(s.connectionInfo, s.dbName, s.db)
	s.server.Close()
	s.connection.Close()
}

func TestIntegrationSuite(test *testing.T) {
	suite.Run(test, new(IntegrationTestSuite))
}
