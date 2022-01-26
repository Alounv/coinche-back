package app

type GameServiceType interface {
	ListGames() []Game
	GetGame(id int) Game
	CreateGame(name string) int
}

type Game struct {
	Name string
	Id   int
}
