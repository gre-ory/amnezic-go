package model

import "fmt"

// //////////////////////////////////////////////////
// error

var (
	ErrNotImplemented          = fmt.Errorf("not implemented")
	ErrGameNotFound            = fmt.Errorf("game not found")
	ErrConcurrentUpdate        = fmt.Errorf("concurrent update")
	ErrInvalidGameId           = fmt.Errorf("invalid game id")
	ErrInvalidMusicId          = fmt.Errorf("invalid music id")
	ErrInvalidThemeId          = fmt.Errorf("invalid theme id")
	ErrInvalidDeezerId         = fmt.Errorf("invalid deezer id")
	ErrInvalidNbPlayer         = fmt.Errorf("invalid number of player")
	ErrInvalidNbQuestion       = fmt.Errorf("invalid number of question")
	ErrInvalidNbAnswer         = fmt.Errorf("invalid number of answer")
	ErrMissingSource           = fmt.Errorf("missing source")
	ErrMusicNotFound           = fmt.Errorf("music not found")
	ErrMusicAlbumNotFound      = fmt.Errorf("music album not found")
	ErrMusicArtistNotFound     = fmt.Errorf("music artist not found")
	ErrMusicGenreNotFound      = fmt.Errorf("music genre not found")
	ErrThemeNotFound           = fmt.Errorf("theme not found")
	ErrThemeQuestionNotFound   = fmt.Errorf("theme question not found")
	ErrInvalidBody             = fmt.Errorf("invalid body")
	ErrInvalidNumberOfQuestion = fmt.Errorf("invalid number of question")
	ErrInvalidNumberOfAnswer   = fmt.Errorf("invalid number of answer")
	ErrMusicUsed               = fmt.Errorf("music used")
)
