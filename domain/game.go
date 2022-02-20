package domain

import (
	"errors"
	"time"
)

type Phase int

const (
	Preparation Phase = 0
	Teaming     Phase = 1
	Bidding     Phase = 2
	Playing     Phase = 3
	Counting    Phase = 4
	Pause       Phase = 5
)

type Value int

const (
	Eight    Value = 80
	Nine     Value = 90
	Ten      Value = 100
	Eleven   Value = 110
	Twelve   Value = 120
	Thirteen Value = 130
	Fourteen Value = 140
	Fifteen  Value = 150
	Capot    Value = 160
)

type Color string

const (
	Club     Color = "club"
	Diamond  Color = "diamond"
	Heart    Color = "heart"
	Spade    Color = "spade"
	NoTrump  Color = "noTrump"
	AllTrump Color = "allTrump"
)

type Bid struct {
	Player string
	Color  Color
}

type Game struct {
	ID        int
	Name      string
	CreatedAt time.Time
	Players   map[string]Player
	Phase     Phase
	Bids      map[Value]Bid
}

type Player struct {
	Team  string
	Order int
}

func (player Player) CanPlay() bool {
	return player.Order == 0
}

func NewGame(name string) Game {
	return Game{
		Name:    name,
		Players: map[string]Player{},
		Phase:   Preparation,
		Bids:    map[Value]Bid{},
	}
}

const (
	ErrAlreadyInGame   = "ALREADY IN GAME"
	ErrEmptyPlayerName = "EMPTY PLAYER NAME"
	ErrGameFull        = "GAME IS FULL"
	ErrPlayerNotFound  = "PLAYER NOT FOUND"
	ErrNotTeaming      = "NOT IN TEAMING PHASE"
	ErrNotBidding      = "NOT IN BIDDING PHASE"
	ErrTeamFull        = "TEAM IS FULL"
	ErrStartGame       = "GAME CANNOT START"
	ErrTeamsNotEqual   = "TEAMS ARE NOT EQUAL"
	ErrBidTooSmall     = "BID IS TOO SMALL"
	ErrNotYourTurn     = "NOT YOUR TURN"
)

func (game Game) IsFull() bool {
	return len(game.Players) == 4
}

func (game Game) CanStart() error {
	if game.Phase != Teaming {
		return errors.New(ErrNotTeaming)
	}

	team1 := ""
	team1Size := 0

	team2 := ""
	team2Size := 0

	for _, player := range game.Players {
		if player.Team == "" {
			continue
		}

		if team1 == "" {
			team1 = player.Team
			team1Size++
			continue
		}

		if team1 == player.Team {
			team1Size++
			continue
		}

		if team2 == "" {
			team2 = player.Team
			team2Size++
			continue
		}

		if team2 == player.Team {
			team2Size++
			continue
		}

		return errors.New(ErrTeamsNotEqual)
	}

	if team1Size == 2 && team2Size == 2 {
		return nil
	}

	return errors.New(ErrTeamsNotEqual)
}

func (game *Game) AddPlayer(playerName string) error {
	if playerName == "" {
		return errors.New(ErrEmptyPlayerName)
	}

	if _, ok := game.Players[playerName]; ok {
		return errors.New(ErrAlreadyInGame)
	}

	if game.IsFull() {
		return errors.New(ErrGameFull)
	}

	game.Players[playerName] = Player{}
	if game.IsFull() && game.Phase == Preparation {
		game.Phase = Teaming
	}
	return nil
}

func (game *Game) RemovePlayer(playerName string) error {
	if _, ok := game.Players[playerName]; !ok {
		return errors.New(ErrPlayerNotFound)
	}

	delete(game.Players, playerName)

	if !game.IsFull() && game.Phase != Preparation {
		game.Phase = Pause
	}
	return nil
}

func (game *Game) AssignTeam(playerName string, teamName string) error {
	if game.Phase != Teaming {
		return errors.New(ErrNotTeaming)
	}

	teamSize := 0
	for _, player := range game.Players {
		if player.Team == teamName {
			teamSize++
		}
	}

	if teamSize >= 2 {
		return errors.New(ErrTeamFull)
	}

	newPlayer := game.Players[playerName]
	newPlayer.Team = teamName

	game.Players[playerName] = newPlayer

	return nil
}

func (game *Game) ClearTeam(playerName string) error {
	if game.Phase != Teaming {
		return errors.New(ErrNotTeaming)
	}

	newPlayer := game.Players[playerName]
	newPlayer.Team = ""

	game.Players[playerName] = newPlayer

	return nil
}

func (game *Game) Start() error {
	err := game.CanStart()
	if err != nil {
		return err
	}

	game.Phase = Bidding

	shouldInitiateOrder := false
	for _, player := range game.Players {
		if player.Order == 0 {
			shouldInitiateOrder = true
		}
		break
	}

	if shouldInitiateOrder {
		game.initiateOrder()
	} else {
		game.rotateOrder()
	}

	return nil
}

func (game *Game) PlaceBid(player string, value Value, color Color) error {
	if game.Phase != Bidding {
		return errors.New(ErrNotBidding)
	}

	var maxValue Value
	for maxValue = range game.Bids {
		break
	}

	if value <= maxValue {
		return errors.New(ErrBidTooSmall)
	}

	err := game.checkPlayerTurn(player)
	if err != nil {
		return err
	}

	game.Bids[value] = Bid{
		Player: player,
		Color:  color,
	}

	game.rotateOrder()

	return nil
}

func (game *Game) rotateOrder() {
	for name, player := range game.Players {
		if player.Order == 4 {
			player.Order = 1
		} else {
			player.Order++
		}

		game.Players[name] = player
	}
}

func (game *Game) initiateOrder() {
	team1 := ""
	team2 := ""
	for name, player := range game.Players {

		if team1 == "" {
			team1 = player.Team
			player.Order = 1
		} else if team1 == player.Team {
			player.Order = 3
		} else if team2 == "" {
			team2 = player.Team
			player.Order = 2
		} else {
			player.Order = 4
		}

		game.Players[name] = player
	}
}

func (game *Game) checkPlayerTurn(playerName string) error {
	if game.Players[playerName].Order != 0 {
		return errors.New(ErrNotYourTurn)
	}
	return nil
}
