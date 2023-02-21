package model

// //////////////////////////////////////////////////
// question

type Question struct {
	Id      QuestionId
	Theme   Theme
	Music   Music
	Answers []*Answer
}
