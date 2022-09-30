package main

import (
	"coinche/api"
	"coinche/domain"
	repository "coinche/repository"
	"coinche/usecases"
	testUtilities "coinche/utilities/test"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	epg "github.com/fergusstrange/embedded-postgres"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const testLogPrefix = "\n\n ****** TEST: "

type IntegrationTestSuite struct {
	suite.Suite
	db           *sqlx.DB
	dbName       string
	router       *gin.Engine
	gameUsecases *usecases.GameUsecases
	server1      *httptest.Server
	server2      *httptest.Server
	server3      *httptest.Server
	server4      *httptest.Server
	connection1  *websocket.Conn
	connection2  *websocket.Conn
	connection3  *websocket.Conn
	connection4  *websocket.Conn
	hub          *api.Hub
	lastTestGame domain.Game
	postgres     *epg.EmbeddedPostgres
}

func TestIntegrationSuite(test *testing.T) {
	suite.Run(test, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.dbName = "testdb"
	s.db, s.postgres = testUtilities.CreateDb(s.dbName)

	gameRepository, err := repository.NewGameRepositoryFromDb(s.db)
	if err != nil {
		s.T().Fatal(err)
	}

	s.gameUsecases = &usecases.GameUsecases{Repo: gameRepository}

	s.router, s.hub = api.SetupRouter(s.gameUsecases)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	testUtilities.DropDb(s.postgres, s.dbName, s.db)
	s.server1.Close()
	s.server2.Close()
	s.server3.Close()
	s.server4.Close()
	s.connection1.Close()
	s.connection2.Close()
	s.connection3.Close()
	s.connection4.Close()

	err := s.postgres.Stop()
	fmt.Println(err)
}

func (s *IntegrationTestSuite) TestCreateGame() {
	test := s.T()
	assert := assert.New(test)
	response := httptest.NewRecorder()

	test.Run("create game", func(test *testing.T) {
		fmt.Println(testLogPrefix, "create game")
		s.router.ServeHTTP(httptest.NewRecorder(), testUtilities.NewCreateGameRequest(test, "NEW GAME"))

		assert.Equal(http.StatusOK, response.Code)
	})

	test.Run("get game", func(test *testing.T) {
		fmt.Println(testLogPrefix, "get game")
		s.router.ServeHTTP(response, testUtilities.NewGetGameRequest(test, 1))

		got := testUtilities.DecodeToGame(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal("NEW GAME", got.Name)
		assert.Equal(1, got.ID)
		assert.Equal(map[string]domain.Player{}, got.Players)
		assert.IsType([]domain.CardID{}, got.Deck)
		assert.IsType(time.Time{}, got.CreatedAt)
	})

	test.Run("list game", func(test *testing.T) {
		fmt.Println(testLogPrefix, "list game")
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
		fmt.Println(testLogPrefix, "join game")
		s.server1, s.connection1 = api.NewGameWebSocketServer(test, 1, "P1", s.hub)
		receivedGame := api.ReceiveGameOrFatal(s.connection1, test)

		assert.IsType(domain.Game{}, receivedGame)

		assert.Equal(1, len(receivedGame.Players))
		assert.Equal(map[string]domain.Player{"P1": {Hand: []domain.CardID{}}}, receivedGame.Players)
	})

	test.Run("leave unstarted game", func(test *testing.T) {
		fmt.Println(testLogPrefix, "leave unstarted game")
		err := api.SendMessage(s.connection1, "leave", "P1")
		if err != nil {
			test.Fatal(err)
		}

		message := api.ReceiveMessageOrFatal(s.connection1, test)

		assert.Equal("P1 has left the game", message)

		s.router.ServeHTTP(response, testUtilities.NewGetGameRequest(test, 1))
		got := testUtilities.DecodeToGame(response.Body, test)

		assert.Equal(0, len(got.Players))
	})

	test.Run("join back unstarted game", func(test *testing.T) {
		fmt.Println(testLogPrefix, "join back unstarted game")
		s.server1, s.connection1 = api.NewGameWebSocketServer(test, 1, "P1", s.hub)
		got := api.ReceiveGameOrFatal(s.connection1, test)

		assert.IsType(domain.Game{}, got)

		assert.Equal(1, len(got.Players))
		assert.Equal(map[string]domain.Player{"P1": {Hand: []domain.CardID{}}}, got.Players)
	})

	test.Run("other players join", func(test *testing.T) {
		fmt.Println(testLogPrefix, "other players join")
		s.server2, s.connection2 = api.NewGameWebSocketServer(test, 1, "P2", s.hub)

		time.Sleep(50 * time.Millisecond) // wait because if all players join at the same time, the IsFull never gets true in AddPlayer

		s.server3, s.connection3 = api.NewGameWebSocketServer(test, 1, "P3", s.hub)

		time.Sleep(50 * time.Millisecond)

		s.server4, s.connection4 = api.NewGameWebSocketServer(test, 1, "P4", s.hub)

		api.ReceiveMultipleGameOrFatal(s.connection1, test, 3)
		api.ReceiveMultipleGameOrFatal(s.connection2, test, 3)
		api.ReceiveMultipleGameOrFatal(s.connection3, test, 2)
		got := api.ReceiveGameOrFatal(s.connection4, test)

		assert.IsType(domain.Game{}, got)

		assert.Equal(4, len(got.Players))
		assert.Equal(0, len(got.Players["P1"].Hand))
		assert.Equal(0, len(got.Players["P2"].Hand))
		assert.Equal(0, len(got.Players["P3"].Hand))
		assert.Equal(0, len(got.Players["P4"].Hand))
		assert.Equal(domain.Teaming, got.Phase)
	})

	test.Run("start game should return error if team not ready", func(test *testing.T) {
		fmt.Println(testLogPrefix, "start game should return error if team not ready")
		err := api.SendMessage(s.connection1, "start", "P1")
		if err != nil {
			test.Fatal(err)
		}

		got := api.ReceiveMessageOrFatal(s.connection1, test)

		assert.Equal("Could not start game: TEAMS ARE NOT EQUAL", got)
	})

	test.Run("join team", func(test *testing.T) {
		fmt.Println(testLogPrefix, "join game")
		err := api.SendMessage(s.connection1, "joinTeam: Odd", "P1")
		if err != nil {
			test.Fatal(err)
		}
		err = api.SendMessage(s.connection2, "joinTeam: Even", "P2")
		if err != nil {
			test.Fatal(err)
		}
		err = api.SendMessage(s.connection3, "joinTeam: Odd", "P3")
		if err != nil {
			test.Fatal(err)
		}
		err = api.SendMessage(s.connection4, "joinTeam: Even", "P4")
		if err != nil {
			test.Fatal(err)
		}

		api.ReceiveMultipleGameOrFatal(s.connection1, test, 3)
		api.ReceiveMultipleGameOrFatal(s.connection2, test, 4)
		api.ReceiveMultipleGameOrFatal(s.connection3, test, 4)
		api.ReceiveMultipleGameOrFatal(s.connection4, test, 4)
		got := api.ReceiveGameOrFatal(s.connection1, test)

		assert.Equal(map[string]domain.Player{"P1": {Team: "Odd", Hand: []domain.CardID{}}, "P2": {Team: "Even", Hand: []domain.CardID{}}, "P3": {Team: "Odd", Hand: []domain.CardID{}}, "P4": {Team: "Even", Hand: []domain.CardID{}}}, got.Players)
	})

	test.Run("start game", func(test *testing.T) {
		fmt.Println(testLogPrefix, "start game")
		err := api.SendMessage(s.connection1, "start", "P1")
		if err != nil {
			test.Fatal(err)
		}

		got := api.ReceiveGameOrFatal(s.connection1, test)
		api.ReceiveGameOrFatal(s.connection2, test)
		api.ReceiveGameOrFatal(s.connection3, test)
		api.ReceiveGameOrFatal(s.connection4, test)

		assert.Equal(1, got.ID)
		assert.Equal("NEW GAME", got.Name)
		assert.Equal(domain.Bidding, got.Phase)
		assert.Equal(map[domain.BidValue]domain.Bid{}, got.Bids)
		assert.Equal(0, len(got.Deck))
		assert.Equal(8, len(got.Players["P1"].Hand))
		assert.Equal(1, got.Players["P1"].Order)
		assert.Equal(2, got.Players["P2"].Order)
		assert.Equal(3, got.Players["P3"].Order)
		assert.Equal(4, got.Players["P4"].Order)
	})

	test.Run("place some bids", func(test *testing.T) {
		fmt.Println(testLogPrefix, "place some bids")
		err := api.SendMessage(s.connection1, "bid: spade,80", "P1")
		if err != nil {
			test.Fatal(err)
		}

		got := api.ReceiveGameOrFatal(s.connection4, test)
		api.ReceiveGameOrFatal(s.connection1, test)
		api.ReceiveGameOrFatal(s.connection2, test)
		api.ReceiveGameOrFatal(s.connection3, test)

		assert.Equal(1, got.ID)
		assert.Equal(1, len(got.Bids))
		assert.Equal(0, got.Bids[domain.Eighty].Coinche)
		assert.Equal(0, got.Bids[domain.Eighty].Pass)
		assert.Equal(domain.Spade, got.Bids[domain.Eighty].Color)
		assert.Equal("P1", got.Bids[domain.Eighty].Player)
	})

	test.Run("place some bids with error", func(test *testing.T) {
		fmt.Println(testLogPrefix, "place some bids with error")
		err := api.SendMessage(s.connection2, "bid: heart,80", "P2")
		if err != nil {
			test.Fatal(err)
		}

		got := api.ReceiveMessageOrFatal(s.connection2, test)

		assert.Equal("Could not bid: BID IS TOO SMALL", got)
	})

	test.Run("can start playing", func(test *testing.T) {
		fmt.Println(testLogPrefix, "can start playing")
		err := api.SendMessage(s.connection2, "bid: pass", "P2")
		if err != nil {
			test.Fatal(err)
		}

		time.Sleep(50 * time.Millisecond) // wait to prevent submitting bid at the same time

		err = api.SendMessage(s.connection3, "bid: spade,90", "P3")
		if err != nil {
			test.Fatal(err)
		}

		time.Sleep(50 * time.Millisecond)

		err = api.SendMessage(s.connection2, "bid: coinche", "P2")
		if err != nil {
			test.Fatal(err)
		}

		time.Sleep(50 * time.Millisecond)

		err = api.SendMessage(s.connection3, "bid: pass", "P3")
		if err != nil {
			test.Fatal(err)
		}

		time.Sleep(50 * time.Millisecond)

		err = api.SendMessage(s.connection1, "bid: pass", "P1")
		if err != nil {
			test.Fatal(err)
		}

		api.ReceiveMultipleGameOrFatal(s.connection1, test, 5)
		api.ReceiveMultipleGameOrFatal(s.connection2, test, 5)
		api.ReceiveMultipleGameOrFatal(s.connection3, test, 5)
		api.ReceiveMultipleGameOrFatal(s.connection4, test, 4)

		got := api.ReceiveGameOrFatal(s.connection4, test)

		assert.Equal(1, got.ID)
		assert.Equal(2, len(got.Bids))
		assert.Equal(1, got.Bids[domain.Ninety].Coinche)
		assert.Equal(2, got.Bids[domain.Ninety].Pass)
		assert.Equal(domain.Spade, got.Bids[domain.Ninety].Color)
		assert.Equal("P3", got.Bids[domain.Ninety].Player)

		assert.Equal(domain.Playing, got.Phase)
		assert.Equal(1, got.Players["P1"].Order)
		assert.Equal(8, len(got.Players["P1"].Hand))

	})

	test.Run("other players are notified when a player leaves", func(test *testing.T) {
		fmt.Println(testLogPrefix, "other players are notified when a player leaves")
		// SHOULD WORK ALSO WITH 	s.connection1.Close()

		err := api.SendMessage(s.connection1, "leave", "P1")
		if err != nil {
			test.Fatal(err)
		}

		message := api.ReceiveMessageOrFatal(s.connection1, test)
		assert.Equal("P1 has left the game", message)

		api.ReceiveMessageOrFatal(s.connection2, test)
		api.ReceiveMessageOrFatal(s.connection3, test)
		api.ReceiveMessageOrFatal(s.connection4, test)

		api.ReceiveGameOrFatal(s.connection1, test)
		api.ReceiveGameOrFatal(s.connection2, test)
		api.ReceiveGameOrFatal(s.connection3, test)
		api.ReceiveGameOrFatal(s.connection4, test)
	})

	test.Run("player can go back in the game", func(test *testing.T) {
		fmt.Println(testLogPrefix, "player can go back in the game")
		s.server1, s.connection1 = api.NewGameWebSocketServer(test, 1, "P1", s.hub)

		api.ReceiveGameOrFatal(s.connection2, test)
		api.ReceiveGameOrFatal(s.connection3, test)
		api.ReceiveGameOrFatal(s.connection4, test)
		got := api.ReceiveGameOrFatal(s.connection1, test)

		assert.Equal(1, got.ID)
		assert.Equal(2, len(got.Bids))
		assert.Equal(1, got.Bids[domain.Ninety].Coinche)
		assert.Equal(2, got.Bids[domain.Ninety].Pass)
		assert.Equal(domain.Spade, got.Bids[domain.Ninety].Color)
		assert.Equal("P3", got.Bids[domain.Ninety].Player)

		assert.Equal(domain.Playing, got.Phase)
		assert.Equal(1, got.Players["P1"].Order)
		assert.Equal(8, len(got.Players["P1"].Hand))

		s.lastTestGame = got
	})
}

func (s *IntegrationTestSuite) TestPlayGame() {
	test := s.T()
	assert := assert.New(test)
	game := s.lastTestGame

	test.Run("can play cards", func(test *testing.T) {
		fmt.Println(testLogPrefix, "can play cards")
		playerHand := game.Players["P1"].Hand
		card := string(playerHand[0])

		err := api.SendMessage(s.connection1, fmt.Sprint("play: ", card), "P1")
		if err != nil {
			test.Fatal(err)
		}

		api.ReceiveGameOrFatal(s.connection1, test)
		api.ReceiveGameOrFatal(s.connection2, test)
		api.ReceiveGameOrFatal(s.connection3, test)

		got := api.ReceiveGameOrFatal(s.connection4, test)

		assert.Equal(1, len(got.Turns))

		assert.Equal(7, len(got.Players["P1"].Hand))
		assert.Equal(8, len(got.Players["P2"].Hand))
		assert.Equal(8, len(got.Players["P3"].Hand))
		assert.Equal(8, len(got.Players["P4"].Hand))

		s.lastTestGame = got
	})

	test.Run("can play all cards", func(test *testing.T) {
		fmt.Println(testLogPrefix, "can play all cards")
		game := s.lastTestGame
		connections := []*websocket.Conn{s.connection1, s.connection2, s.connection3, s.connection4}

		for t := 0; t < 8; t++ {
			fmt.Println("\n TURN ", t)
			sortedPlayerNames := testUtilities.GetSortedPlayersNameByOrder(game.Players)

			for _, playerName := range sortedPlayerNames {
				p := testUtilities.GetPlayerIndexFromNameOrFatal(playerName, test)
				playerHand := game.Players[playerName].Hand

				for c := 0; c < len(playerHand); c++ {
					card := string(playerHand[c])

					err := api.SendMessage(connections[p], fmt.Sprint("play: ", card), playerName)
					if err != nil {
						test.Fatal(err)
					}

					message, newGame := api.ReceiveMessageOrGameOrFatal(connections[p], test)

					if message == "" { // We did not receive an error message, so we can update the game
						api.EmtpyGamesForOtherPlayersOrFatal(sortedPlayerNames, playerName, 1, test, connections)
						game = newGame
						break
					}
				}
			}
		}

		assert.Equal(0, len(game.Players["P1"].Hand))
		assert.Equal(0, len(game.Players["P2"].Hand))
		assert.Equal(0, len(game.Players["P3"].Hand))
		assert.Equal(0, len(game.Players["P4"].Hand))

		assert.Equal(8, len(game.Turns))
		assert.Equal(domain.Counting, game.Phase)

		s.lastTestGame = game
	})

	test.Run("can count points", func(test *testing.T) {
		fmt.Println(testLogPrefix, "can count points")
		game := s.lastTestGame

		totalPoints := game.Points["Odd"] + game.Points["Even"]
		totalScores := game.Scores["Odd"] + game.Scores["Even"]

		if totalPoints == 162 {
			assert.Equal(true, totalScores == 500 || totalScores == 540)
		} else if totalPoints == 182 {
			assert.Equal(540, totalScores)
		} else {
			test.Fatal("Points and scores are not equal to 162 or 182")
		}
	})
}

// TODO: TEST RESTART

// TODO: TEST DISCONNECTION IN GAME
