package model

import "fmt"

// //////////////////////////////////////////////////
// genre

type Genre struct {
	Id     int64
	Name   string
	ImgUrl string
}

func (obj *Genre) String() string {
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
