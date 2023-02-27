package model

import "fmt"

// //////////////////////////////////////////////////
// answer

type Answer struct {
	Id      AnswerId
	Text    string
	Hint    string
	Correct bool
}

func (obj *Answer) String() string {
	if obj == nil {
		return ""
	}
	return fmt.Sprintf(
		"{ \"id\": %d, \"text\": \"%s\", \"hint\": \"%s\", \"correct\": %t }",
		obj.Id,
		obj.Text,
		obj.Hint,
		obj.Correct,
	)
}
