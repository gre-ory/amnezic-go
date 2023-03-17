package model

// //////////////////////////////////////////////////
// constants

const (
	GameIdStart     = 10000000
	QuestionIdStart = 10000
	AnswerIdStart   = 100
	PlayerIdStart   = 1
)

// //////////////////////////////////////////////////
// game id

type GameId int64

func NewGameId(number int) GameId {
	return GameId(GameIdStart * number)
}

// //////////////////////////////////////////////////
// question id

type QuestionId int64

func NewQuestionId(gameId GameId, number int) QuestionId {
	return QuestionId(int64(gameId) + int64(QuestionIdStart*number))
}

func (id QuestionId) Split() GameId {
	gameId := GameId((id / GameIdStart) * GameIdStart)
	return gameId
}

// //////////////////////////////////////////////////
// answer id

type AnswerId int64

func NewAnswerId(questionId QuestionId, number int) AnswerId {
	return AnswerId(int64(questionId) + int64(AnswerIdStart*number))
}

func (id AnswerId) Split() (GameId, QuestionId) {
	gameId := GameId((id / GameIdStart) * GameIdStart)
	questionId := QuestionId((id / QuestionIdStart) * QuestionIdStart)
	return gameId, questionId
}

// //////////////////////////////////////////////////
// player answer id

type PlayerAnswerId int64

func NewPlayerAnswerId(answerId AnswerId, playerId PlayerId) PlayerAnswerId {
	return PlayerAnswerId(int64(answerId) + int64(playerId))
}

func (id PlayerAnswerId) Split() (GameId, QuestionId, AnswerId, PlayerId) {
	gameId := GameId((id / GameIdStart) * GameIdStart)
	questionId := QuestionId((id / QuestionIdStart) * QuestionIdStart)
	answerId := AnswerId((id / AnswerIdStart) * AnswerIdStart)
	playerId := PlayerId(id % AnswerIdStart)
	return gameId, questionId, answerId, playerId
}

// //////////////////////////////////////////////////
// player id

type PlayerId int64

func NewPlayerId(number int) PlayerId {
	return PlayerId(PlayerIdStart * number)
}
