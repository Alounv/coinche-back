package domain

import (
	"errors"
)

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

func (game *Game) AddPlayer(playerName string) error {
	if game.Phase != Teaming {
		return errors.New(ErrNotTeaming)
	}

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

	if game.IsFull() && game.Phase == Teaming {
		game.Phase = Teaming
	}
	return nil
}

func (game *Game) RemovePlayer(playerName string) error {
	if game.Phase != Teaming {
		return errors.New(ErrNotTeaming)
	}

	if _, ok := game.Players[playerName]; !ok {
		return errors.New(ErrPlayerNotFound)
	}

	delete(game.Players, playerName)

	return nil
}

func (game Game) CanStartBidding() error {
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

	if game.CanStartBidding() == nil {
		game.Deck = NewDeck()
	}

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
