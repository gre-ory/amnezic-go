package model

// //////////////////////////////////////////////////
// game

type Game struct {
	Id        GameId
	Version   int
	Settings  *GameSettings
	Players   []*Player
	Questions []*Question
}
