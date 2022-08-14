package domain

import (
	"errors"
	"math/rand"
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

type BidValue int

const (
	Eighty           BidValue = 80
	Ninety           BidValue = 90
	Hundred          BidValue = 100
	HundredAndTen    BidValue = 110
	HundredAndTwenty BidValue = 120
	HundredAndThirty BidValue = 130
	HundredAndFourty BidValue = 140
	HundredAndFifty  BidValue = 150
	Capot            BidValue = 160
)

type Strength int

const (
	Seven  Strength = 1
	Eight  Strength = 2
	Nine   Strength = 3
	Jack   Strength = 4
	Queen  Strength = 5
	King   Strength = 6
	Ten    Strength = 7
	As     Strength = 8
	TSeven Strength = 11
	TEight Strength = 12
	TQueen Strength = 13
	TKing  Strength = 14
	TTen   Strength = 15
	TAs    Strength = 16
	TNine  Strength = 17
	TJack  Strength = 18
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

type card struct {
	color         Color
	strength      Strength
	TrumpStrength Strength
	value         int
	trumpValue    int
}

func (card card) getValue(trump Color) int {
	if trump == card.color || trump == AllTrump {
		return card.trumpValue
	} else {
		return card.value
	}
}

type CardID string

const (
	C_7  CardID = "7-club"
	C_8  CardID = "8-club"
	C_9  CardID = "9-club"
	C_10 CardID = "10-club"
	C_J  CardID = "jack-club"
	C_Q  CardID = "queen-club"
	C_K  CardID = "king-club"
	C_A  CardID = "as-club"
	D_7  CardID = "7-diamond"
	D_8  CardID = "8-diamond"
	D_9  CardID = "9-diamond"
	D_10 CardID = "10-diamond"
	D_J  CardID = "jack-diamond"
	D_Q  CardID = "queen-diamond"
	D_K  CardID = "king-diamond"
	D_A  CardID = "as-diamond"
	H_7  CardID = "7-heart"
	H_8  CardID = "8-heart"
	H_9  CardID = "9-heart"
	H_10 CardID = "10-heart"
	H_J  CardID = "jack-heart"
	H_Q  CardID = "queen-heart"
	H_K  CardID = "king-heart"
	H_A  CardID = "as-heart"
	S_7  CardID = "7-spade"
	S_8  CardID = "8-spade"
	S_9  CardID = "9-spade"
	S_10 CardID = "10-spade"
	S_J  CardID = "jack-spade"
	S_Q  CardID = "queen-spade"
	S_K  CardID = "king-spade"
	S_A  CardID = "as-spade"
)

var cards = map[CardID]card{
	C_7:  {Club, Seven, TSeven, 0, 0},
	C_8:  {Club, Eight, TEight, 0, 0},
	C_9:  {Club, Nine, TNine, 0, 14},
	C_10: {Club, Ten, TTen, 10, 10},
	C_J:  {Club, Jack, TJack, 2, 20},
	C_Q:  {Club, Queen, TQueen, 3, 3},
	C_K:  {Club, King, TKing, 4, 4},
	C_A:  {Club, As, TAs, 11, 11},
	D_7:  {Diamond, Seven, TSeven, 0, 0},
	D_8:  {Diamond, Eight, TEight, 0, 0},
	D_9:  {Diamond, Nine, TNine, 0, 14},
	D_10: {Diamond, Ten, TTen, 10, 10},
	D_J:  {Diamond, Jack, TJack, 2, 20},
	D_Q:  {Diamond, Queen, TQueen, 3, 3},
	D_K:  {Diamond, King, TKing, 4, 4},
	D_A:  {Diamond, As, TAs, 11, 11},
	H_7:  {Heart, Seven, TSeven, 0, 0},
	H_8:  {Heart, Eight, TEight, 0, 0},
	H_9:  {Heart, Nine, TNine, 0, 14},
	H_10: {Heart, Ten, TTen, 10, 10},
	H_J:  {Heart, Jack, TJack, 2, 20},
	H_Q:  {Heart, Queen, TQueen, 3, 3},
	H_K:  {Heart, King, TKing, 4, 4},
	H_A:  {Heart, As, TAs, 11, 11},
	S_7:  {Spade, Seven, TSeven, 0, 0},
	S_8:  {Spade, Eight, TEight, 0, 0},
	S_9:  {Spade, Nine, TNine, 0, 14},
	S_10: {Spade, Ten, TTen, 10, 10},
	S_J:  {Spade, Jack, TJack, 2, 20},
	S_Q:  {Spade, Queen, TQueen, 3, 3},
	S_K:  {Spade, King, TKing, 4, 4},
	S_A:  {Spade, As, TAs, 11, 11},
}

func newDeck() []CardID {
	deck := []CardID{C_7, C_8, C_9, C_10, C_J, C_Q, C_K, C_A, D_7, D_8, D_9, D_10, D_J, D_Q, D_K, D_A, H_7, H_8, H_9, H_10, H_J, H_Q, H_K, H_A, S_7, S_8, S_9, S_10, S_J, S_Q, S_K, S_A}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

type Bid struct {
	Player  string
	Color   Color
	Coinche int
	Pass    int
}

type play struct {
	playerName string
	card       CardID
}

type turn struct {
	plays  []play
	winner string
}

func (turn turn) getWinner(trump Color) string {
	var winner string
	var strongerValue Strength
	var firstCard CardID
	for _, play := range turn.plays {
		if firstCard == "" {
			firstCard = play.card
		}
		cardValue := getCardValue(play.card, trump, firstCard)
		if cardValue > strongerValue {
			strongerValue = cardValue
			winner = play.playerName
		}
	}

	return winner
}

func (turn *turn) setWinner(trump Color) {
	turn.winner = turn.getWinner(trump)
}

func getCardValue(card CardID, trump Color, firstCard CardID) Strength {
	color := cards[card].color
	colorAsked := cards[firstCard].color

	if trump == color || trump == AllTrump {
		return cards[card].TrumpStrength
	} else if color == colorAsked {
		return cards[card].strength
	} else {
		return 0
	}
}

type Game struct {
	ID        int
	Name      string
	CreatedAt time.Time
	Players   map[string]Player
	Phase     Phase
	Bids      map[BidValue]Bid
	Deck      []CardID
	turns     []turn
	scores    map[string]int
	points    map[string]int
}

type Player struct {
	Team         string
	Order        int
	InitialOrder int
	Hand         []CardID
}

func (player Player) CanPlay() bool {
	return player.Order == 0
}

func NewGame(name string) Game {
	return Game{
		Name:    name,
		Players: map[string]Player{},
		Phase:   Preparation,
		Bids:    map[BidValue]Bid{},
		Deck:    newDeck(),
	}
}

func (game *Game) checkPlayerTurn(playerName string) error {
	if game.Players[playerName].Order != 1 {
		return errors.New(ErrNotYourTurn)
	}
	return nil
}

func (game *Game) checkTeamTurn(playerName string) error {
	order := game.Players[playerName].Order
	if order != 1 && order != 3 {
		return errors.New(ErrNotYourTurn)
	}
	return nil
}

func (game *Game) setFirstPlayer(playerName string) {
	for i := 0; i < len(game.Players); i++ {
		if game.Players[playerName].Order == 1 {
			return
		}
		game.rotateOrder()
	}
}

func (game *Game) rotateOrder() {
	for name, player := range game.Players {
		if player.Order == 1 {
			player.Order = 4
		} else {
			player.Order--
		}

		game.Players[name] = player
	}
}

func (game *Game) trump() Color {
	lastBid, _ := game.getLastBid()
	return lastBid.Color
}
