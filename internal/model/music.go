package model

// //////////////////////////////////////////////////
// music

type Music struct {
	Id     int64
	Name   string
	Mp3Url string
	Artist *Artist
	Album  *Album
	Genre  *Genre
}
