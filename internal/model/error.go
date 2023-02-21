package model

import "fmt"

// //////////////////////////////////////////////////
// error

var (
	ErrNotImplemented    = fmt.Errorf("not implemented")
	ErrGameNotFound      = fmt.Errorf("game not found")
	ErrConcurrentUpdate  = fmt.Errorf("concurrent update")
	ErrInvalidGameId     = fmt.Errorf("invalid game id")
	ErrInvalidNbPlayer   = fmt.Errorf("invalid number of player")
	ErrInvalidNbQuestion = fmt.Errorf("invalid number of question")
	ErrInvalidNbAnswer   = fmt.Errorf("invalid number of answer")
)
