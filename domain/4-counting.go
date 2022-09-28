package domain

const (
	CAPO_WON_SCORE  = 500
	CAPO_LOST_SCORE = 320
)

func (game *Game) end() {
	game.Phase = Counting

	game.calculatesTeamPointsAndScores()
}

func (game Game) getPotentialBelotes() map[Color]string {
	potentialBelotes := map[Color]string{}

	trump := game.trump()
	playersCards := game.getPlayersCards()

	for player, playerCards := range playersCards {
		potentialPlayerBelotes := map[Color]int{}

		for _, card := range playerCards {
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

	return potentialBelotes
}

func (game Game) getPlayerWithBelote() string {
	playerWithBelote := ""

	potentialBelotes := game.getPotentialBelotes()

	if len(potentialBelotes) > 0 {
		for _, turn := range game.Turns {
			for _, play := range turn.Plays {
				card := cards[play.Card]
				cardStrength := card.getStrength(game.trump())
				if cardStrength == TQueen || cardStrength == TKing {
					if player, ok := potentialBelotes[card.color]; ok {
						playerWithBelote = player
					}
				}
			}
		}
	}

	return playerWithBelote
}

func (game Game) getTeams() (string, string) {
	lastBid, _ := game.getLastBid()
	contractTeam := game.Players[lastBid.Player].Team
	otherTeam := ""

	for team := range game.Points {
		if team != contractTeam {
			otherTeam = team
		}
	}

	return contractTeam, otherTeam

}

func (game *Game) calculateBasePoints() {
	trump := game.trump()
	playersCards := game.getPlayersCards()

	for player, playerCards := range playersCards {
		team := game.Players[player].Team

		for _, card := range playerCards {
			game.Points[team] += cards[card].getValue(trump)
		}
	}
}

func (game *Game) applyLastTen() {
	lastTurn := game.Turns[len(game.Turns)-1]
	lastWinnerTeam := game.Players[lastTurn.Winner].Team
	game.Points[lastWinnerTeam] += 10
}

func (game *Game) applyAllTrumpNoTrump(contractTeam string) {
	trump := game.trump()

	if trump == NoTrump {
		game.Points[contractTeam] = game.Points[contractTeam] * 162 / 130 // converting to int automatically rounds down which is what we want because we use >= to check if contract is fulfilled
	} else if trump == AllTrump {
		game.Points[contractTeam] = game.Points[contractTeam] * 162 / 258
	}
}

func (game *Game) calculatesTeamPoints() (contractTeamPointsWithoutBelote int, otherTeamPointsWithoutBelote int) {
	game.Points = map[string]int{}

	game.calculateBasePoints()

	game.applyLastTen()

	contractTeam, otherTeam := game.getTeams()

	game.applyAllTrumpNoTrump(contractTeam)

	game.Points[otherTeam] = 162 - game.Points[contractTeam]

	contractTeamPointsWithoutBelote = game.Points[contractTeam]
	otherTeamPointsWithoutBelote = game.Points[otherTeam]

	game.applyBeloteToPoints(contractTeam)

	return contractTeamPointsWithoutBelote, otherTeamPointsWithoutBelote
}

func (game *Game) applyBeloteToPoints(contractTeam string) {
	playerWithBelote := game.getPlayerWithBelote()

	if playerWithBelote != "" {
		team := game.Players[playerWithBelote].Team

		if contractTeam == team {
			game.Points[team] += 20
		}
	}
}

func (game *Game) applyBeloteToScores() {
	playerWithBelote := game.getPlayerWithBelote()

	if playerWithBelote != "" {
		team := game.Players[playerWithBelote].Team
		game.Scores[team] += 20
	}
}

func (game *Game) addContractPoints(isCapot bool, isContractWon bool, contractPoints int) {
	contractTeam, otherTeam := game.getTeams()

	if !isCapot {
		if isContractWon {
			game.Scores[contractTeam] += contractPoints
		} else {
			game.Scores[otherTeam] += contractPoints
		}
	}
}

func (game *Game) addRealizedPoints(isCapot bool, isContractWon bool, coinche int, contractTeamPointsWithoutBelote int, otherTeamPointsWithoutBelote int) {
	contractTeam, otherTeam := game.getTeams()

	isCoinche := coinche > 0
	isCapotWon := isCapot && game.Points[otherTeam] == 0
	isCapotLost := isCapot && game.Points[otherTeam] != 0
	isNormalContractWon := !isCapot && isContractWon
	isNormalContractLost := !isCapot && !isContractWon
	isNormalContractWonWithCoinche := isNormalContractWon && isCoinche
	isNormalContractLostWithCoinche := isNormalContractLost && isCoinche

	if isCapotWon {
		game.Scores[contractTeam] += CAPO_WON_SCORE
	} else if isCapotLost {
		game.Scores[otherTeam] += CAPO_LOST_SCORE
	} else if isNormalContractWonWithCoinche {
		game.Scores[contractTeam] += 160
	} else if isNormalContractLostWithCoinche {
		game.Scores[otherTeam] += 160
	} else if isNormalContractWon {
		game.Scores[contractTeam] += contractTeamPointsWithoutBelote
		game.Scores[otherTeam] += otherTeamPointsWithoutBelote
	} else {
		game.Scores[contractTeam] += contractTeamPointsWithoutBelote
		game.Scores[otherTeam] += 160
	}
}

func (game *Game) applyCoinche(coinche int) {
	for team, score := range game.Scores {
		game.Scores[team] = getScoreWithCoinche(score, coinche)
	}
}

func (game *Game) calculatesTeamScores(contractTeamPointsWithoutBelote int, otherTeamPointsWithoutBelote int) {
	contractTeam, otherTeam := game.getTeams()
	game.Scores = map[string]int{
		contractTeam: 0,
		otherTeam:    0,
	}

	lastBid, contract := game.getLastBid()
	contractPoints := int(contract)
	isCapot := contract == Capot
	isContractWon := game.Points[contractTeam] >= contractPoints

	game.addContractPoints(isCapot, isContractWon, contractPoints)

	game.applyBeloteToScores()

	coinche := lastBid.Coinche

	game.addRealizedPoints(isCapot, isContractWon, coinche, contractTeamPointsWithoutBelote, otherTeamPointsWithoutBelote)

	game.applyCoinche(coinche)
}

func (game *Game) calculatesTeamPointsAndScores() {
	contractTeamPointsWithoutBelote, otherTeamPointsWithoutBelote := game.calculatesTeamPoints()
	game.calculatesTeamScores(contractTeamPointsWithoutBelote, otherTeamPointsWithoutBelote)
}

func (card card) getStrength(trump Color) Strength {
	if trump == card.color || trump == AllTrump {
		return card.TrumpStrength
	}

	return card.strength
}

func (game Game) getPlayersCards() map[string][]CardID {
	playersCards := map[string][]CardID{}

	for _, turn := range game.Turns {
		for _, play := range turn.Plays {
			playersCards[turn.Winner] = append(playersCards[turn.Winner], play.Card)
		}
	}

	return playersCards
}

func getScoreWithCoinche(score int, coinche int) int {
	if coinche > 0 {
		return score * 2 * coinche
	}

	return score
}
