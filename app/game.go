package app

type GameServiceType interface {
	GetAllGames() []Game
	GetAGame(id int) Game
	CreateAGame(name string) int
}

type Game struct {
	Name  string
	Id int
}