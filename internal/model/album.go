package model

import "fmt"

// //////////////////////////////////////////////////
// album

type Album struct {
	Id     int64
	Name   string
	ImgUrl string
}

func (obj *Album) String() string {
	if obj == nil {
		return ""
	}
	return fmt.Sprintf(
		"{ \"id\": %d, \"name\": \"%s\", \"img-url\": \"%s\" }",
		obj.Id,
		obj.Name,
		obj.ImgUrl,
	)
}
