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
// game question id

type GameQuestionId int64

func NewGameQuestionId(gameId GameId, number int) GameQuestionId {
	return GameQuestionId(int64(gameId) + int64(QuestionIdStart*number))
}

func (id GameQuestionId) Split() GameId {
	gameId := GameId((id / GameIdStart) * GameIdStart)
	return gameId
}

// //////////////////////////////////////////////////
// game answer id

type GameAnswerId int64

func NewGameAnswerId(questionId GameQuestionId, number int) GameAnswerId {
	return GameAnswerId(int64(questionId) + int64(AnswerIdStart*number))
}

func (id GameAnswerId) Split() (GameId, GameQuestionId) {
	gameId := GameId((id / GameIdStart) * GameIdStart)
	questionId := GameQuestionId((id / QuestionIdStart) * QuestionIdStart)
	return gameId, questionId
}

// //////////////////////////////////////////////////
// player answer id

type GamePlayerAnswerId int64

func NewGamePlayerAnswerId(answerId GameAnswerId, playerId GamePlayerId) GamePlayerAnswerId {
	return GamePlayerAnswerId(int64(answerId) + int64(playerId))
}

func (id GamePlayerAnswerId) Split() (GameId, GameQuestionId, GameAnswerId, GamePlayerId) {
	gameId := GameId((id / GameIdStart) * GameIdStart)
	questionId := GameQuestionId((id / QuestionIdStart) * QuestionIdStart)
	answerId := GameAnswerId((id / AnswerIdStart) * AnswerIdStart)
	playerId := GamePlayerId(id % AnswerIdStart)
	return gameId, questionId, answerId, playerId
}

// //////////////////////////////////////////////////
// game player id

type GamePlayerId int64

func NewGamePlayerId(number int) GamePlayerId {
	return GamePlayerId(PlayerIdStart * number)
}
