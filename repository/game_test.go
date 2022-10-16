package repository

import (
	"coinche/domain"
	"coinche/utilities"
	testUtilities "coinche/utilities/test"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func NewGameRepositoryWithData(db *sqlx.DB) (*GameRepository, error) {
	repository, err := NewGameRepositoryFromDb(db)
	if err != nil {
		return nil, err
	}

	games := []domain.Game{
		{Name: "GAME ONE", ID: 1, Players: map[string]domain.Player{}},
		{Name: "GAME TWO", ID: 2, Players: map[string]domain.Player{"P1": {}, "P2": {}}},
		newTeamingGame(),
		newCompleteGame(),
		{Name: "PREVIOUS GAME ONE", ID: 3, Players: map[string]domain.Player{}, Root: 1},
	}

	tx := repository.db.MustBegin()

	for _, game := range games {
		_, err := createGame(game, tx)
		if err != nil {
			return nil, err
		}
	}

	return repository, tx.Commit()
}

func newTeamingGame() domain.Game {
	return domain.Game{
		ID:      3,
		Name:    "GAME TEAMING",
		Players: map[string]domain.Player{},
		Phase:   domain.Teaming,
		Bids:    map[domain.BidValue]domain.Bid{},
		Deck: []domain.CardID{
			domain.C7, domain.C8, domain.C9, domain.C10, domain.CJ, domain.CQ, domain.CK, domain.CA,
			domain.D7, domain.D8, domain.D9, domain.D10, domain.DJ, domain.DQ, domain.DK, domain.DA,
			domain.H7, domain.H8, domain.H9, domain.H10, domain.HJ, domain.HQ, domain.HK, domain.HA,
			domain.S7, domain.S8, domain.S9, domain.S10, domain.SJ, domain.SQ, domain.SK, domain.SA,
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
			{PlayerName: "P1", Card: domain.C7},
			{PlayerName: "P2", Card: domain.C10},
			{PlayerName: "P3", Card: domain.CK},
			{PlayerName: "P4", Card: domain.H9},
		}, Winner: "P4"},
		{Plays: []domain.Play{
			{PlayerName: "P4", Card: domain.D8},
			{PlayerName: "P1", Card: domain.DJ},
			{PlayerName: "P2", Card: domain.DK},
			{PlayerName: "P3", Card: domain.D7},
		}, Winner: "P2"},
		{Plays: []domain.Play{
			{PlayerName: "P2", Card: domain.CJ},
			{PlayerName: "P3", Card: domain.CA},
			{PlayerName: "P4", Card: domain.H10},
			{PlayerName: "P1", Card: domain.C8},
		}, Winner: "P4"},
		{Plays: []domain.Play{
			{PlayerName: "P4", Card: domain.D9},
			{PlayerName: "P1", Card: domain.DQ},
			{PlayerName: "P2", Card: domain.DA},
			{PlayerName: "P3", Card: domain.H7},
		}, Winner: "P3"},
		{Plays: []domain.Play{
			{PlayerName: "P3", Card: domain.H8},
			{PlayerName: "P4", Card: domain.D10},
			{PlayerName: "P1", Card: domain.HJ},
			{PlayerName: "P2", Card: domain.HA},
		}, Winner: "P1"},
		{Plays: []domain.Play{
			{PlayerName: "P1", Card: domain.C9},
			{PlayerName: "P2", Card: domain.CQ},
			{PlayerName: "P3", Card: domain.S9},
			{PlayerName: "P4", Card: domain.SQ},
		}, Winner: "P2"},
		{Plays: []domain.Play{
			{PlayerName: "P2", Card: domain.S7},
			{PlayerName: "P3", Card: domain.S10},
			{PlayerName: "P4", Card: domain.SK},
			{PlayerName: "P1", Card: domain.HQ},
		}, Winner: "P1"},
		{Plays: []domain.Play{
			{PlayerName: "P1", Card: domain.S8},
			{PlayerName: "P2", Card: domain.HK},
			{PlayerName: "P3", Card: domain.SJ},
			{PlayerName: "P4", Card: domain.SA},
		}, Winner: "P2"},
	}
	return game
}

func TestGameRepo(test *testing.T) {
	assert := assert.New(test)
	dbName := "testgamerepodb"
	utilities.LoadEnv("../.env")

	db, postgres := testUtilities.CreateDb(dbName)

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
		testUtilities.DropDb(postgres, dbName, db)
	})
}

func TestGameRepoWithInitialData(test *testing.T) {
	assert := assert.New(test)
	dbName := "testgamerepowithinitialdatadb"
	db, postgres := testUtilities.CreateDb(dbName)

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

	test.Run("list all games, where root equal id", func(test *testing.T) {
		want := []domain.Game{
			{Name: "GAME ONE", ID: 1, Players: map[string]domain.Player{}, Root: 1},
			{Name: "GAME TWO", ID: 2, Players: map[string]domain.Player{"P1": {Team: "", Order: 0, InitialOrder: 0, Hand: []domain.CardID(nil)}, "P2": {Team: "", Order: 0, InitialOrder: 0, Hand: []domain.CardID(nil)}}, Root: 2},
		}

		got, err := repository.ListGames()
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(4, len(got))
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

	test.Run("update a game", func(test *testing.T) {
		want := domain.Game{
			ID:    2,
			Phase: domain.Bidding,
			Bids: map[domain.BidValue]domain.Bid{
				domain.Eighty: {Player: "P1", Color: domain.Spade, Coinche: 1},
			},
			Players: map[string]domain.Player{
				"P1": {Hand: []domain.CardID{domain.C7}, Order: 1, InitialOrder: 1, Team: "A Team"},
				"P2": {Hand: []domain.CardID{}},
				"P3": {Hand: []domain.CardID{}},
				"P4": {Hand: []domain.CardID{}},
			},
			Turns: []domain.Turn{
				{Plays: []domain.Play{
					{PlayerName: "P1", Card: domain.C7},
					{PlayerName: "P2", Card: domain.C10},
				}, Winner: "P4"},
			},
			Points: map[string]int{
				"A Team": 80,
				"B Team": 82,
			},
			Scores: map[string]int{
				"A Team": 1000,
				"B Team": 500,
			},
			Root: 2,
		}

		err := repository.UpdateGame(want)
		if err != nil {
			test.Fatal(err)
		}

		want.Turns = []domain.Turn{
			{Plays: []domain.Play{
				{PlayerName: "P1", Card: domain.C7},
				{PlayerName: "P2", Card: domain.C10},
				{PlayerName: "P3", Card: domain.CK},
				{PlayerName: "P4", Card: domain.H9},
			}, Winner: "P4"},
			{Plays: []domain.Play{
				{PlayerName: "P4", Card: domain.D8},
				{PlayerName: "P1", Card: domain.DJ},
				{PlayerName: "P2", Card: domain.DK},
				{PlayerName: "P3", Card: domain.D7},
			}, Winner: "P2"},
		}

		want.Bids = map[domain.BidValue]domain.Bid{
			domain.Eighty: {Player: "P1", Color: domain.Spade, Coinche: 2, Pass: 3},
		}

		want.Root = 0

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
		assert.Equal(want.Points, got.Points)
		assert.Equal(want.Scores, got.Scores)
		assert.Equal(want.Root, got.Root)
	})

	test.Run("reset a game", func(test *testing.T) {
		want := domain.Game{
			ID:     2,
			Phase:  domain.Bidding,
			Bids:   map[domain.BidValue]domain.Bid{},
			Turns:  []domain.Turn{},
			Points: map[string]int{},
		}

		err := repository.UpdateGame(want)
		if err != nil {
			test.Fatal(err)
		}

		got, err := repository.GetGame(2)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(0, len(got.Bids))
		assert.Equal(0, len(got.Turns))
		assert.Equal(0, len(got.Points))
		assert.Equal(0, len(got.Deck))
	})

	test.Cleanup(func() {
		testUtilities.DropDb(postgres, dbName, db)
	})
}
