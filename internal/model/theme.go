package model

import "fmt"

// //////////////////////////////////////////////////
// theme

type Theme struct {
	Id    int64
	Title string
}

func (obj *Theme) String() string {
	if obj == nil {
		return ""
	}
	return fmt.Sprintf(
		"{ \"id\": %d, \"title\": \"%s\" }",
		obj.Id,
		obj.Title,
	)
}
