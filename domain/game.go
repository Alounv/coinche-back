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
}

var cards = map[int]card{
	0:  {Club, Seven, TSeven},
	1:  {Club, Eight, TEight},
	2:  {Club, Nine, TNine},
	3:  {Club, Ten, TTen},
	4:  {Club, Jack, TJack},
	5:  {Club, Queen, TQueen},
	6:  {Club, King, TKing},
	7:  {Club, As, TAs},
	8:  {Diamond, Seven, TSeven},
	9:  {Diamond, Eight, TEight},
	10: {Diamond, Nine, TNine},
	11: {Diamond, Ten, TTen},
	12: {Diamond, Jack, TJack},
	13: {Diamond, Queen, TQueen},
	14: {Diamond, King, TKing},
	15: {Diamond, As, TAs},
	16: {Heart, Seven, TSeven},
	17: {Heart, Eight, TEight},
	18: {Heart, Nine, TNine},
	19: {Heart, Ten, TTen},
	20: {Heart, Jack, TJack},
	21: {Heart, Queen, TQueen},
	22: {Heart, King, TKing},
	23: {Heart, As, TAs},
	24: {Spade, Seven, TSeven},
	25: {Spade, Eight, TEight},
	26: {Spade, Nine, TNine},
	27: {Spade, Ten, TTen},
	28: {Spade, Jack, TJack},
	29: {Spade, Queen, TQueen},
	30: {Spade, King, TKing},
	31: {Spade, As, TAs},
}

func shuffle() []int {
	a := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a
}

type Bid struct {
	Player  string
	Color   Color
	Coinche int
	Pass    int
}

type play struct {
	playerName string
	card       int
}

type turn struct {
	plays  []play
	winner string
}

func (turn *turn) setWinner(trump Color) {
	var winner string
	var strongerValue Strength
	for _, play := range turn.plays {
		cardValue := getCardValue(play.card, trump)
		if cardValue > strongerValue {
			strongerValue = cardValue
			winner = play.playerName
		}
	}

	turn.winner = winner
}

func getCardValue(card int, trump Color) Strength {
	color := cards[card].color

	if trump == color || trump == AllTrump {
		return cards[card].TrumpStrength
	} else {
		return cards[card].strength
	}
}

type Game struct {
	ID        int
	Name      string
	CreatedAt time.Time
	Players   map[string]Player
	Phase     Phase
	Bids      map[BidValue]Bid
	trump     Color
	deck      []int
	turns     []turn
}

type Player struct {
	Team         string
	Order        int
	InitialOrder int
	Hand         []int
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
		deck:    shuffle(),
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
