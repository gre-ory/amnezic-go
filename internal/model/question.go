package model

import (
	"fmt"

	"github.com/gre-ory/amnezic-go/internal/util"
)

// //////////////////////////////////////////////////
// question

type Question struct {
	Id      QuestionId
	Theme   Theme
	Music   Music
	Answers []*Answer
}

func (obj *Question) String() string {
	if obj == nil {
		return ""
	}
	return fmt.Sprintf(
		"{ \"id\": %d, \"theme\": %s, \"music\": %s, \"answers\": [%s] }",
		obj.Id,
		obj.Theme.String(),
		obj.Music.String(),
		util.Join(obj.Answers, (*Answer).String),
	)
}
