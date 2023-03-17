package model

import "fmt"

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

func (obj *Music) String() string {
	if obj == nil {
		return ""
	}
	return fmt.Sprintf(
		"{ \"id\": %d, \"name\": \"%s\", \"mp3-url\": \"%s\", \"artist\": %s, \"album\": %s, \"genre\": %s }",
		obj.Id,
		obj.Name,
		obj.Mp3Url,
		obj.Artist,
		obj.Album,
		obj.Genre,
	)
}
