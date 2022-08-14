package domain

import (
	"fmt"
)

const (
	CAPO_WON_SCORE  = 500
	CAPO_LOST_SCORE = 320
)

func (game *Game) end() {
	game.Phase = Counting

	game.calculatesTeamPoints()

	fmt.Println(game.points, game.scores)
}

func (game *Game) calculatesTeamPoints() {
	playersCards := game.getPlayersCards()

	game.points = map[string]int{}

	potentialBelotes := map[Color]string{}

	trump := game.trump()

	for player, playerCards := range playersCards {
		team := game.Players[player].Team
		potentialPlayerBelotes := map[Color]int{}

		for _, card := range playerCards {
			game.points[team] += cards[card].getValue(trump)

			card := cards[card]
			cardStrength := card.getStrength(trump)
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
	game.points[lastWinnerTeam] += 10

	lastBid, contract := game.getLastBid()
	contractTeam := game.Players[lastBid.Player].Team
	otherTeam := ""

	for team := range game.points {
		if team != contractTeam {
			otherTeam = team
		}
	}

	if trump == NoTrump {
		game.points[contractTeam] = game.points[contractTeam] * 162 / 130 // converting to int automatically rounds down which is what we want because we use >= to check if contract is fulfilled
	} else if trump == AllTrump {
		game.points[contractTeam] = game.points[contractTeam] * 162 / 258
	}

	game.points[otherTeam] = 162 - game.points[contractTeam]

	// IF PLAYER TEAM AS TAKEN
	playerWithBelote := ""

	if len(potentialBelotes) > 0 {
		for _, turn := range game.turns {
			for _, play := range turn.plays {
				card := cards[play.card]
				cardStrength := card.getStrength(trump)
				if cardStrength == TQueen || cardStrength == TKing {
					if player, ok := potentialBelotes[card.color]; ok {
						playerWithBelote = player
					}
				}
			}
		}
	}

	game.scores = map[string]int{
		contractTeam: 0,
		otherTeam:    0,
	}

	contractTeamPointsWithoutBelote := game.points[contractTeam]
	otherTeamPointsWithoutBelote := game.points[otherTeam]

	if playerWithBelote != "" {
		team := game.Players[playerWithBelote].Team
		game.scores[team] += 20

		if contractTeam == team {
			game.points[team] += 20
		}
	}

	coinche := lastBid.Coinche
	isCapot := contract == Capot
	isCoinche := coinche > 0
	contractPoints := int(contract)

	isCapotWon := isCapot && game.points[otherTeam] == 0
	isCapotLost := isCapot && game.points[otherTeam] != 0
	isContractWon := game.points[contractTeam] >= contractPoints
	isNormalContractWon := !isCapot && isContractWon
	isNormalContractLost := !isCapot && !isContractWon
	isNormalContractWonWithCoinche := isNormalContractWon && isCoinche
	isNormalContractLostWithCoinche := isNormalContractLost && isCoinche

	if !isCapot {
		if isContractWon {
			game.scores[contractTeam] += contractPoints
		} else {
			game.scores[otherTeam] += contractPoints
		}
	}

	if isCapotWon {
		game.scores[contractTeam] += CAPO_WON_SCORE
	} else if isCapotLost {
		game.scores[otherTeam] += CAPO_LOST_SCORE
	} else if isNormalContractWonWithCoinche {
		game.scores[contractTeam] += 160
	} else if isNormalContractLostWithCoinche {
		game.scores[otherTeam] += 160
	} else if isNormalContractWon {
		game.scores[contractTeam] += contractTeamPointsWithoutBelote
		game.scores[otherTeam] += otherTeamPointsWithoutBelote
	} else {
		game.scores[contractTeam] += contractTeamPointsWithoutBelote
		game.scores[otherTeam] += 160
	}

	for team, score := range game.scores {
		game.scores[team] = getScoreWithCoinche(score, coinche)
	}
}

func (card card) getStrength(trump Color) Strength {
	if trump == card.color || trump == AllTrump {
		return card.TrumpStrength
	} else {
		return card.strength
	}
}

func (game Game) getPlayersCards() map[string][]CardID {
	playersCards := map[string][]CardID{}

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
