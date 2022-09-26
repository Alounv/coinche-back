package domain

import (
	"errors"
	"fmt"
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
	C7  CardID = "7-club"
	C8  CardID = "8-club"
	C9  CardID = "9-club"
	C10 CardID = "10-club"
	CJ  CardID = "jack-club"
	CQ  CardID = "queen-club"
	CK  CardID = "king-club"
	CA  CardID = "as-club"
	D7  CardID = "7-diamond"
	D8  CardID = "8-diamond"
	D9  CardID = "9-diamond"
	D10 CardID = "10-diamond"
	DJ  CardID = "jack-diamond"
	DQ  CardID = "queen-diamond"
	DK  CardID = "king-diamond"
	DA  CardID = "as-diamond"
	H7  CardID = "7-heart"
	H8  CardID = "8-heart"
	H9  CardID = "9-heart"
	H10 CardID = "10-heart"
	HJ  CardID = "jack-heart"
	HQ  CardID = "queen-heart"
	HK  CardID = "king-heart"
	HA  CardID = "as-heart"
	S7  CardID = "7-spade"
	S8  CardID = "8-spade"
	S9  CardID = "9-spade"
	S10 CardID = "10-spade"
	SJ  CardID = "jack-spade"
	SQ  CardID = "queen-spade"
	SK  CardID = "king-spade"
	SA  CardID = "as-spade"
)

var cards = map[CardID]card{
	C7:  {Club, Seven, TSeven, 0, 0},
	C8:  {Club, Eight, TEight, 0, 0},
	C9:  {Club, Nine, TNine, 0, 14},
	C10: {Club, Ten, TTen, 10, 10},
	CJ:  {Club, Jack, TJack, 2, 20},
	CQ:  {Club, Queen, TQueen, 3, 3},
	CK:  {Club, King, TKing, 4, 4},
	CA:  {Club, As, TAs, 11, 11},
	D7:  {Diamond, Seven, TSeven, 0, 0},
	D8:  {Diamond, Eight, TEight, 0, 0},
	D9:  {Diamond, Nine, TNine, 0, 14},
	D10: {Diamond, Ten, TTen, 10, 10},
	DJ:  {Diamond, Jack, TJack, 2, 20},
	DQ:  {Diamond, Queen, TQueen, 3, 3},
	DK:  {Diamond, King, TKing, 4, 4},
	DA:  {Diamond, As, TAs, 11, 11},
	H7:  {Heart, Seven, TSeven, 0, 0},
	H8:  {Heart, Eight, TEight, 0, 0},
	H9:  {Heart, Nine, TNine, 0, 14},
	H10: {Heart, Ten, TTen, 10, 10},
	HJ:  {Heart, Jack, TJack, 2, 20},
	HQ:  {Heart, Queen, TQueen, 3, 3},
	HK:  {Heart, King, TKing, 4, 4},
	HA:  {Heart, As, TAs, 11, 11},
	S7:  {Spade, Seven, TSeven, 0, 0},
	S8:  {Spade, Eight, TEight, 0, 0},
	S9:  {Spade, Nine, TNine, 0, 14},
	S10: {Spade, Ten, TTen, 10, 10},
	SJ:  {Spade, Jack, TJack, 2, 20},
	SQ:  {Spade, Queen, TQueen, 3, 3},
	SK:  {Spade, King, TKing, 4, 4},
	SA:  {Spade, As, TAs, 11, 11},
}

func NewDeck() []CardID {
	deck := []CardID{C7, C8, C9, C10, CJ, CQ, CK, CA, D7, D8, D9, D10, DJ, DQ, DK, DA, H7, H8, H9, H10, HJ, HQ, HK, HA, S7, S8, S9, S10, SJ, SQ, SK, SA}
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

type Play struct {
	PlayerName string
	Card       CardID
}

type Turn struct {
	Plays  []Play
	Winner string
}

func (turn Turn) getWinner(trump Color) string {
	var winner string
	var strongerValue Strength
	var firstCard CardID
	for _, play := range turn.Plays {
		if firstCard == "" {
			firstCard = play.Card
		}
		cardValue := getCardValue(play.Card, trump, firstCard)
		if cardValue > strongerValue {
			strongerValue = cardValue
			winner = play.PlayerName
		}
	}

	return winner
}

func (turn *Turn) setWinner(trump Color) {
	turn.Winner = turn.getWinner(trump)
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
	Turns     []Turn
	Scores    map[string]int
	Points    map[string]int
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
		Deck:    NewDeck(),
	}
}

func (game *Game) checkPlayerTurn(playerName string) error {
	order := game.Players[playerName].Order
	if order != 1 {
		return errors.New(fmt.Sprint(ErrNotYourTurn, " ", playerName, " ", order))
	}
	return nil
}

func (game *Game) checkTeamTurn(playerName string) error {
	order := game.Players[playerName].Order
	if order != 1 && order != 3 {
		return errors.New(ErrNotYourTeamTurn)
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
