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

	require.Equal(t, model.GameId(10000000), model.NewGameId(1))
	require.Equal(t, model.GameId(20000000), model.NewGameId(2))
	require.Equal(t, model.GameId(990000000), model.NewGameId(99))
	require.Equal(t, model.GameId(9990000000), model.NewGameId(999))

	gameId := model.NewGameId(42)
	require.Equal(t, model.QuestionId(420010000), model.NewQuestionId(gameId, 1))
	require.Equal(t, model.QuestionId(420020000), model.NewQuestionId(gameId, 2))
	require.Equal(t, model.QuestionId(420990000), model.NewQuestionId(gameId, 99))

	questionId := model.NewQuestionId(gameId, 17)
	require.Equal(t, model.AnswerId(420170100), model.NewAnswerId(questionId, 1))
	require.Equal(t, model.AnswerId(420170200), model.NewAnswerId(questionId, 2))
	require.Equal(t, model.AnswerId(420170900), model.NewAnswerId(questionId, 9))

	answerId := model.NewAnswerId(questionId, 3)
	require.Equal(t, model.PlayerAnswerId(420170301), model.NewPlayerAnswerId(answerId, model.NewPlayerId(1)))
	require.Equal(t, model.PlayerAnswerId(420170302), model.NewPlayerAnswerId(answerId, model.NewPlayerId(2)))
	require.Equal(t, model.PlayerAnswerId(420170309), model.NewPlayerAnswerId(answerId, model.NewPlayerId(9)))

	playerAnswerId := model.NewPlayerAnswerId(answerId, model.NewPlayerId(2))
	gameId, questionId, answerId, playerId := playerAnswerId.Split()
	require.Equal(t, model.GameId(420000000), gameId)
	require.Equal(t, model.QuestionId(420170000), questionId)
	require.Equal(t, model.AnswerId(420170300), answerId)
	require.Equal(t, model.PlayerId(2), playerId)

	gameId, questionId = answerId.Split()
	require.Equal(t, model.GameId(420000000), gameId)
	require.Equal(t, model.QuestionId(420170000), questionId)

	gameId = questionId.Split()
	require.Equal(t, model.GameId(420000000), gameId)
}
