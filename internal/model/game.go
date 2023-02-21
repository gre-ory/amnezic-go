package model

// //////////////////////////////////////////////////
// game

type Game struct {
	Id        GameId
	Version   int
	Players   []*Player
	Questions []*Question
}
