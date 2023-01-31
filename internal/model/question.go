package model

// //////////////////////////////////////////////////
// question

type Question struct {
	Id     int64
	Theme  Theme
	Music  *Music
	Answer []*Answer
}
