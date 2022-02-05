package main

import (
	"coinche/api"
	gameapi "coinche/api/game"
	"coinche/domain"
	gamerepo "coinche/repository/game"
	"coinche/usecases"
	"coinche/utilities/env"
	testutils "coinche/utilities/test"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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
}

func (s *IntegrationTestSuite) SetupSuite() {
	env.LoadEnv("")
	s.connectionInfo = os.Getenv("SQLX_POSTGRES_INFO")
	s.dbName = "testdb"

	s.db = testutils.CreateDb(s.connectionInfo, s.dbName)

	gameRepository := gamerepo.NewGameRepositoryFromDb(s.db)
	s.gameUsecases = &usecases.GameUsecases{Repo: gameRepository}

	s.router = api.SetupRouter(s.gameUsecases)
}

func (s *IntegrationTestSuite) TestCreateGame() {
	test := s.T()
	assert := assert.New(test)

	s.router.ServeHTTP(httptest.NewRecorder(), testutils.NewCreateGameRequest(test, "NEW GAME"))

	response := httptest.NewRecorder()
	assert.Equal(http.StatusOK, response.Code)
}

func (s *IntegrationTestSuite) TestGetGame() {
	test := s.T()
	assert := assert.New(test)

	response := httptest.NewRecorder()
	s.router.ServeHTTP(response, testutils.NewGetGameRequest(test, 1))

	got := testutils.DecodeToGame(response.Body, test)

	assert.Equal(http.StatusOK, response.Code)
	assert.Equal("NEW GAME", got.Name)
	assert.Equal(1, got.ID)
	assert.IsType(time.Time{}, got.CreatedAt)
}

func (s *IntegrationTestSuite) TestListGames() {
	test := s.T()
	assert := assert.New(test)
	request, _ := http.NewRequest(http.MethodGet, "/games/all", nil)
	response := httptest.NewRecorder()

	s.router.ServeHTTP(response, request)

	got := testutils.DecodeToGames(response.Body, test)

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

	server, connection := testutils.NewGameWebSocketServer(s.T(), s.gameUsecases, 1, "player")
	receivedGame, _ := gameapi.ReceiveGame(connection)

	assert.IsType(domain.Game{}, receivedGame)

	s.router.ServeHTTP(response, testutils.NewGetGameRequest(test, 1))
	got := testutils.DecodeToGame(response.Body, test)

	assert.Equal([]string{"player"}, got.Players)

	test.Cleanup(func() {
		server.Close()
		connection.Close()
	})
}

func (s *IntegrationTestSuite) TearDownSuite() {
	testutils.DropDb(s.connectionInfo, s.dbName, s.db)
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
