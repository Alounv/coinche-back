package domain

import (
	"fmt"
)

func (game *Game) end() {
	game.Phase = Counting

	teamsPoints, teamScore := game.getTeamPoints()

	fmt.Println(teamsPoints, teamScore)
}

func (game Game) getTeamPoints() (points map[string]int, scores map[string]int) {
	playersCards := game.getPlayersCards()

	teamPoints := map[string]int{}

	potentialBelotes := map[Color]string{}

	for player, playerCards := range playersCards {
		team := game.Players[player].Team
		potentialPlayerBelotes := map[Color]int{}

		for _, card := range playerCards {
			teamPoints[team] += cards[card].getValue(game.trump)

			card := cards[card]
			cardStrength := card.getStrength(game.trump)
			if cardStrength == TQueen || cardStrength == TKing {
				potentialPlayerBelotes[card.color]++
			}
		}

		for color, number := range potentialPlayerBelotes {
			if number == 2 {
				potentialBelotes[color] = player
			}
		}
	}

	lastTurn := game.turns[len(game.turns)-1]
	lastWinnerTeam := game.Players[lastTurn.winner].Team
	teamPoints[lastWinnerTeam] += 10

	lastBid, contract := game.getLastBid()
	contractTeam := game.Players[lastBid.Player].Team
	otherTeam := ""

	for team := range teamPoints {
		if team != contractTeam {
			otherTeam = team
		}
	}

	if game.trump == NoTrump {
		teamPoints[contractTeam] = teamPoints[contractTeam] * 162 / 130 // converting to int automatically rounds down which is what we want because we use >= to check if contract is fulfilled
	} else if game.trump == AllTrump {
		teamPoints[contractTeam] = teamPoints[contractTeam] * 162 / 258
	}

	teamPoints[otherTeam] = 162 - teamPoints[contractTeam]

	// IF PLAYER TEAM AS TAKEN
	playerWithBelote := ""

	if len(potentialBelotes) > 0 {
		for _, turn := range game.turns {
			for _, play := range turn.plays {
				card := cards[play.card]
				cardStrength := card.getStrength(game.trump)
				if cardStrength == TQueen || cardStrength == TKing {
					if player, ok := potentialBelotes[card.color]; ok {
						playerWithBelote = player
					}
				}
			}
		}
	}

	teamScore := map[string]int{
		contractTeam: 0,
		otherTeam:    0,
	}

	if playerWithBelote != "" {
		team := game.Players[playerWithBelote].Team
		teamScore[team] += 20

		if contractTeam == team {
			teamPoints[team] += 20
		}
	}

	coinche := lastBid.Coinche

	isCapotWon := contract == Capot && teamPoints[otherTeam] == 0
	isNormalContractWon := contract != Capot && teamPoints[contractTeam] >= int(contract)

	if isCapotWon {
		teamScore[contractTeam] += getScoreWithCoinche(160, coinche)
	} else if isNormalContractWon {
		teamScore[contractTeam] += getScoreWithCoinche(int(contract), coinche)
	} else {
		teamScore[otherTeam] += getScoreWithCoinche(160, coinche)
	}

	return teamPoints, teamScore
}

func (card card) getStrength(trump Color) Strength {
	if trump == card.color || trump == AllTrump {
		return card.TrumpStrength
	} else {
		return card.strength
	}
}

func (game Game) getPlayersCards() map[string][]cardID {
	playersCards := map[string][]cardID{}

	for _, turn := range game.turns {
		for _, play := range turn.plays {
			playersCards[turn.winner] = append(playersCards[turn.winner], play.card)
		}
	}

	return playersCards
}

func getScoreWithCoinche(score int, coinche int) int {
	if coinche > 0 {
		return score * 2 * coinche
	} else {
		return score
	}
}
