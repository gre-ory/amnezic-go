package model

// //////////////////////////////////////////////////
// game

type Game struct {
	Id        int64
	Players   []*Player
	Questions []*Question
}
