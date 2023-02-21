package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gre-ory/amnezic-go/internal/model"
)

func TestId(t *testing.T) {

	require.Equal(t, model.PlayerId(1), model.NewPlayerId(1))
	require.Equal(t, model.PlayerId(2), model.NewPlayerId(2))
	require.Equal(t, model.PlayerId(9), model.NewPlayerId(9))

	require.Equal(t, model.GameId(10000), model.NewGameId(1))
	require.Equal(t, model.GameId(20000), model.NewGameId(2))
	require.Equal(t, model.GameId(990000), model.NewGameId(99))

	gameId := model.NewGameId(42)
	require.Equal(t, model.QuestionId(420100), model.NewQuestionId(gameId, 1))
	require.Equal(t, model.QuestionId(420200), model.NewQuestionId(gameId, 2))
	require.Equal(t, model.QuestionId(429900), model.NewQuestionId(gameId, 99))

	questionId := model.NewQuestionId(gameId, 17)
	require.Equal(t, model.AnswerId(421710), model.NewAnswerId(questionId, 1))
	require.Equal(t, model.AnswerId(421720), model.NewAnswerId(questionId, 2))
	require.Equal(t, model.AnswerId(421790), model.NewAnswerId(questionId, 9))

	answerId := model.NewAnswerId(questionId, 3)
	require.Equal(t, model.PlayerAnswerId(421731), model.NewPlayerAnswerId(answerId, model.NewPlayerId(1)))
	require.Equal(t, model.PlayerAnswerId(421732), model.NewPlayerAnswerId(answerId, model.NewPlayerId(2)))
	require.Equal(t, model.PlayerAnswerId(421739), model.NewPlayerAnswerId(answerId, model.NewPlayerId(9)))

	playerAnswerId := model.NewPlayerAnswerId(answerId, model.NewPlayerId(2))
	gameId, questionId, answerId, playerId := playerAnswerId.Split()
	require.Equal(t, model.GameId(420000), gameId)
	require.Equal(t, model.QuestionId(421700), questionId)
	require.Equal(t, model.AnswerId(421730), answerId)
	require.Equal(t, model.PlayerId(2), playerId)

	gameId, questionId = answerId.Split()
	require.Equal(t, model.GameId(420000), gameId)
	require.Equal(t, model.QuestionId(421700), questionId)

	gameId = questionId.Split()
	require.Equal(t, model.GameId(420000), gameId)
}
