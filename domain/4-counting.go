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

	teamsPoints, teamScore := game.getTeamPoints()

	fmt.Println(teamsPoints, teamScore)
}

func (game Game) getTeamPoints() (points map[string]int, scores map[string]int) {
	playersCards := game.getPlayersCards()

	teamPoints := map[string]int{}

	potentialBelotes := map[Color]string{}

	trump := game.trump()

	for player, playerCards := range playersCards {
		team := game.Players[player].Team
		potentialPlayerBelotes := map[Color]int{}

		for _, card := range playerCards {
			teamPoints[team] += cards[card].getValue(trump)

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
	teamPoints[lastWinnerTeam] += 10

	lastBid, contract := game.getLastBid()
	contractTeam := game.Players[lastBid.Player].Team
	otherTeam := ""

	for team := range teamPoints {
		if team != contractTeam {
			otherTeam = team
		}
	}

	if trump == NoTrump {
		teamPoints[contractTeam] = teamPoints[contractTeam] * 162 / 130 // converting to int automatically rounds down which is what we want because we use >= to check if contract is fulfilled
	} else if trump == AllTrump {
		teamPoints[contractTeam] = teamPoints[contractTeam] * 162 / 258
	}

	teamPoints[otherTeam] = 162 - teamPoints[contractTeam]

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

	teamScore := map[string]int{
		contractTeam: 0,
		otherTeam:    0,
	}

	contractTeamPointsWithoutBelote := teamPoints[contractTeam]
	otherTeamPointsWithoutBelote := teamPoints[otherTeam]

	if playerWithBelote != "" {
		team := game.Players[playerWithBelote].Team
		teamScore[team] += 20

		if contractTeam == team {
			teamPoints[team] += 20
		}
	}

	coinche := lastBid.Coinche
	isCapot := contract == Capot
	isCoinche := coinche > 0
	contractPoints := int(contract)

	isCapotWon := isCapot && teamPoints[otherTeam] == 0
	isCapotLost := isCapot && teamPoints[otherTeam] != 0
	isContractWon := teamPoints[contractTeam] >= contractPoints
	isNormalContractWon := !isCapot && isContractWon
	isNormalContractLost := !isCapot && !isContractWon
	isNormalContractWonWithCoinche := isNormalContractWon && isCoinche
	isNormalContractLostWithCoinche := isNormalContractLost && isCoinche

	if !isCapot {
		if isContractWon {
			teamScore[contractTeam] += contractPoints
		} else {
			teamScore[otherTeam] += contractPoints
		}
	}

	if isCapotWon {
		teamScore[contractTeam] += CAPO_WON_SCORE
	} else if isCapotLost {
		teamScore[otherTeam] += CAPO_LOST_SCORE
	} else if isNormalContractWonWithCoinche {
		teamScore[contractTeam] += 160
	} else if isNormalContractLostWithCoinche {
		teamScore[otherTeam] += 160
	} else if isNormalContractWon {
		teamScore[contractTeam] += contractTeamPointsWithoutBelote
		teamScore[otherTeam] += otherTeamPointsWithoutBelote
	} else {
		teamScore[contractTeam] += contractTeamPointsWithoutBelote
		teamScore[otherTeam] += 160
	}

	for team, score := range teamScore {
		teamScore[team] = getScoreWithCoinche(score, coinche)
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
