package repository

import (
	"coinche/domain"
	"coinche/utilities"
	testUtilities "coinche/utilities/test"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newTeamingGame() domain.Game {
	return domain.Game{
		ID:      3,
		Name:    "GAME TEAMING",
		Players: map[string]domain.Player{},
		Phase:   domain.Teaming,
		Bids:    map[domain.BidValue]domain.Bid{},
		Deck: []domain.CardID{
			domain.C_7, domain.C_8, domain.C_9, domain.C_10, domain.C_J, domain.C_Q, domain.C_K, domain.C_A,
			domain.D_7, domain.D_8, domain.D_9, domain.D_10, domain.D_J, domain.D_Q, domain.D_K, domain.D_A,
			domain.H_7, domain.H_8, domain.H_9, domain.H_10, domain.H_J, domain.H_Q, domain.H_K, domain.H_A,
			domain.S_7, domain.S_8, domain.S_9, domain.S_10, domain.S_J, domain.S_Q, domain.S_K, domain.S_A,
		},
	}
}

func newCompleteGame() domain.Game {
	game := newTeamingGame()
	game.ID = 4
	game.Name = "GAME COMPLETE"
	game.Players = map[string]domain.Player{
		"P1": {Team: "odd", Order: 1, InitialOrder: 1},
		"P2": {Team: "even", Order: 2, InitialOrder: 2},
		"P3": {Team: "odd", Order: 3, InitialOrder: 3},
		"P4": {Team: "even", Order: 4, InitialOrder: 4},
	}
	game.Phase = domain.Bidding
	game.Bids = map[domain.BidValue]domain.Bid{
		domain.Eighty: {
			Player:  "P1",
			Color:   domain.Heart,
			Coinche: 0,
			Pass:    0,
		},
	}
	game.Deck = []domain.CardID{}
	game.Scores = map[string]int{
		"odd":  72,
		"even": 90,
	}
	game.Points = map[string]int{
		"odd":  72,
		"even": 160 + 80,
	}
	game.Turns = []domain.Turn{
		{Plays: []domain.Play{
			{PlayerName: "P1", Card: domain.C_7},
			{PlayerName: "P2", Card: domain.C_10},
			{PlayerName: "P3", Card: domain.C_K},
			{PlayerName: "P4", Card: domain.H_9},
		}, Winner: "P4"},
		{Plays: []domain.Play{
			{PlayerName: "P4", Card: domain.D_8},
			{PlayerName: "P1", Card: domain.D_J},
			{PlayerName: "P2", Card: domain.D_K},
			{PlayerName: "P3", Card: domain.D_7},
		}, Winner: "P2"},
		{Plays: []domain.Play{
			{PlayerName: "P2", Card: domain.C_J},
			{PlayerName: "P3", Card: domain.C_A},
			{PlayerName: "P4", Card: domain.H_10},
			{PlayerName: "P1", Card: domain.C_8},
		}, Winner: "P4"},
		{Plays: []domain.Play{
			{PlayerName: "P4", Card: domain.D_9},
			{PlayerName: "P1", Card: domain.D_Q},
			{PlayerName: "P2", Card: domain.D_A},
			{PlayerName: "P3", Card: domain.H_7},
		}, Winner: "P3"},
		{Plays: []domain.Play{
			{PlayerName: "P3", Card: domain.H_8},
			{PlayerName: "P4", Card: domain.D_10},
			{PlayerName: "P1", Card: domain.H_J},
			{PlayerName: "P2", Card: domain.H_A},
		}, Winner: "P1"},
		{Plays: []domain.Play{
			{PlayerName: "P1", Card: domain.C_9},
			{PlayerName: "P2", Card: domain.C_Q},
			{PlayerName: "P3", Card: domain.S_9},
			{PlayerName: "P4", Card: domain.S_Q},
		}, Winner: "P2"},
		{Plays: []domain.Play{
			{PlayerName: "P2", Card: domain.S_7},
			{PlayerName: "P3", Card: domain.S_10},
			{PlayerName: "P4", Card: domain.S_K},
			{PlayerName: "P1", Card: domain.H_Q},
		}, Winner: "P1"},
		{Plays: []domain.Play{
			{PlayerName: "P1", Card: domain.S_8},
			{PlayerName: "P2", Card: domain.H_K},
			{PlayerName: "P3", Card: domain.S_J},
			{PlayerName: "P4", Card: domain.S_A},
		}, Winner: "P2"},
	}
	return game
}

func TestGameRepo(test *testing.T) {
	assert := assert.New(test)
	dbName := "testgamerepodb"
	utilities.LoadEnv("../.env")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")

	db := testUtilities.CreateDb(connectionInfo, dbName)

	repository, err := NewGameRepositoryFromDb(db)
	if err != nil {
		test.Fatal(err)
	}

	test.Run("create a simple game", func(test *testing.T) {
		newName := "NEW GAME ONE"
		newPlayers := map[string]domain.Player{"P1": {}, "P2": {}}

		newID, err := repository.CreateGame(domain.Game{Name: newName, Players: newPlayers})
		if err != nil {
			test.Fatal(err)
		}

		got, err := repository.GetGame(newID)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(newName, got.Name)
		assert.Equal(newPlayers, got.Players)
		assert.Equal(newID, got.ID)
		assert.Equal(domain.Phase(0), got.Phase)
		assert.IsType(time.Time{}, got.CreatedAt)
	})

	test.Run("create a teaming game", func(test *testing.T) {
		newGame := newTeamingGame()

		newID, err := repository.CreateGame(newGame)
		if err != nil {
			test.Fatal(err)
		}

		got, err := repository.GetGame(newID)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(newID, got.ID)
		assert.Equal(newGame.Name, got.Name)
		assert.Equal(newGame.Players, got.Players)
		assert.IsType(time.Time{}, got.CreatedAt)
		assert.Equal(newGame.Phase, got.Phase)
		assert.Equal(newGame.Deck, got.Deck)
	})

	test.Run("create a complete game", func(test *testing.T) {
		newGame := newCompleteGame()

		newID, err := repository.CreateGame(newGame)
		if err != nil {
			test.Fatal(err)
		}

		got, err := repository.GetGame(newID)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(newID, got.ID)
		assert.Equal(newGame.Name, got.Name)
		assert.Equal(newGame.Players, got.Players)
		assert.IsType(time.Time{}, got.CreatedAt)
		assert.Equal(newGame.Phase, got.Phase)
		assert.Equal(newGame.Deck, got.Deck)
		assert.Equal(newGame.Bids, got.Bids)
		assert.Equal(newGame.Turns, got.Turns)
		assert.Equal(newGame.Points, got.Points)
		assert.Equal(newGame.Scores, got.Scores)
	})

	test.Cleanup(func() {
		testUtilities.DropDb(connectionInfo, dbName, db)
	})
}

func TestGameRepoWithInitialData(test *testing.T) {
	assert := assert.New(test)
	dbName := "testgamerepowithinitialdatadb"
	utilities.LoadEnv("../.env")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")

	db := testUtilities.CreateDb(connectionInfo, dbName)

	repository, err := NewGameRepositoryWithData(db)
	if err != nil {
		test.Fatal(err)
	}

	test.Run("get an empty game", func(test *testing.T) {
		want := domain.Game{Name: "GAME ONE", ID: 1, Players: map[string]domain.Player{}}

		got, err := repository.GetGame(1)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(want.ID, got.ID)
		assert.Equal(want.Name, got.Name)
		assert.Equal(want.Players, got.Players)
	})

	test.Run("get a game", func(test *testing.T) {
		want := domain.Game{Name: "GAME TWO", ID: 2, Players: map[string]domain.Player{"P1": {Hand: []domain.CardID(nil)}, "P2": {Hand: []domain.CardID(nil)}}}

		got, err := repository.GetGame(2)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(want.ID, got.ID)
		assert.Equal(want.Name, got.Name)
		assert.Equal(want.Players, got.Players)
	})

	test.Run("list all games", func(test *testing.T) {
		want := []domain.Game{
			{Name: "GAME ONE", ID: 1, Players: map[string]domain.Player{}},
			{Name: "GAME TWO", ID: 2, Players: map[string]domain.Player{"P1": {Team: "", Order: 0, InitialOrder: 0, Hand: []domain.CardID(nil)}, "P2": {Team: "", Order: 0, InitialOrder: 0, Hand: []domain.CardID(nil)}}},
		}

		got, err := repository.ListGames()
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(want[0].ID, got[0].ID)
		assert.Equal(want[0].Players, got[0].Players)
		assert.Equal(want[1].ID, got[1].ID)
		assert.Equal(want[1].Players, got[1].Players)
	})

	test.Run("update a player", func(test *testing.T) {
		player := domain.Player{Team: "A Team"}

		err := repository.UpdatePlayer(2, "P2", player)
		if err != nil {
			test.Fatal(err)
		}

		game, err := repository.GetGame(2)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal("A Team", game.Players["P2"].Team)
	})

	test.Run("update a game", func(test *testing.T) { // FIXME: should be progressively augmented
		want := domain.Game{
			ID:    2,
			Phase: domain.Bidding,
			Bids: map[domain.BidValue]domain.Bid{
				domain.Eighty: {Player: "P1", Color: domain.Spade, Coinche: 1},
			},
			Players: map[string]domain.Player{
				"P1": {Hand: []domain.CardID{domain.C_7}, Order: 1, InitialOrder: 1, Team: "A Team"},
				"P2": {Hand: []domain.CardID{}},
				"P3": {Hand: []domain.CardID{}},
				"P4": {Hand: []domain.CardID{}},
			},
			Turns: []domain.Turn{
				{Plays: []domain.Play{
					{PlayerName: "P1", Card: domain.C_7},
					{PlayerName: "P2", Card: domain.C_10},
				}, Winner: "P4"},
			},
		}

		err := repository.UpdateGame(want)
		if err != nil {
			test.Fatal(err)
		}

		want.Turns = []domain.Turn{
			{Plays: []domain.Play{
				{PlayerName: "P1", Card: domain.C_7},
				{PlayerName: "P2", Card: domain.C_10},
				{PlayerName: "P3", Card: domain.C_K},
				{PlayerName: "P4", Card: domain.H_9},
			}, Winner: "P4"},
			{Plays: []domain.Play{
				{PlayerName: "P4", Card: domain.D_8},
				{PlayerName: "P1", Card: domain.D_J},
				{PlayerName: "P2", Card: domain.D_K},
				{PlayerName: "P3", Card: domain.D_7},
			}, Winner: "P2"},
		}

		err = repository.UpdateGame(want)
		if err != nil {
			test.Fatal(err)
		}

		got, err := repository.GetGame(2)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(want.ID, got.ID)
		assert.Equal(want.Phase, got.Phase)
		assert.Equal(want.Bids, got.Bids)
		assert.Equal(want.Players, got.Players)
		assert.Equal(want.Deck, got.Deck)
		assert.Equal(want.Bids, got.Bids)
		assert.Equal(want.Turns[0], got.Turns[0])
		assert.Equal(want.Turns[1], got.Turns[1])
		/*
			assert.Equal(want.Points, got.Points)
			assert.Equal(want.Scores, got.Scores)
		*/
	})

	test.Cleanup(func() {
		testUtilities.DropDb(connectionInfo, dbName, db)
	})
}

func NewGameRepositoryWithData(db *sqlx.DB) (*GameRepository, error) {
	repository, err := NewGameRepositoryFromDb(db)
	if err != nil {
		return nil, err
	}

	err = repository.CreateGames([]domain.Game{
		{Name: "GAME ONE", ID: 1, Players: map[string]domain.Player{}},
		{Name: "GAME TWO", ID: 2, Players: map[string]domain.Player{"P1": {}, "P2": {}}},
		newTeamingGame(),
		newCompleteGame(),
	})

	return repository, err
}
