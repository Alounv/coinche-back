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

type Game struct {
	ID        int
	Name      string
	CreatedAt time.Time
	Players   map[string]Player
	Phase     Phase
}

type Player struct {
	Team string
}

const (
	ErrAlreadyInGame   = "ALREADY IN GAME"
	ErrEmptyPlayerName = "EMPTY PLAYER NAME"
	ErrGameFull        = "GAME IS FULL"
	ErrPlayerNotFound  = "PLAYER NOT FOUND"
	ErrNotTeaming      = "NOT IN TEAMING PHASE"
	ErrTeamFull        = "TEAM IS FULL"
	ErrStartGame       = "GAME CANNOT START"
	ErrTeamsNotEqual   = "TEAMS ARE NOT EQUAL"
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
		} else if team1 == "" {
			team1 = player.Team
			team1Size++
		} else if team1 == player.Team {
			team1Size++
		} else if team2 == "" {
			team2 = player.Team
			team2Size++
		} else if team2 == player.Team {
			team2Size++
		}
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
	return nil
}

func NewGame(name string) Game {
	return Game{
		Name:    name,
		Players: map[string]Player{},
		Phase:   Preparation,
	}
}
