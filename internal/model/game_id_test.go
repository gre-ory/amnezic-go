package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gre-ory/amnezic-go/internal/model"
)

func TestId(t *testing.T) {

	require.Equal(t, model.GamePlayerId(1), model.NewGamePlayerId(1))
	require.Equal(t, model.GamePlayerId(2), model.NewGamePlayerId(2))
	require.Equal(t, model.GamePlayerId(9), model.NewGamePlayerId(9))

	require.Equal(t, model.GameId(10000000), model.NewGameId(1))
	require.Equal(t, model.GameId(20000000), model.NewGameId(2))
	require.Equal(t, model.GameId(990000000), model.NewGameId(99))
	require.Equal(t, model.GameId(9990000000), model.NewGameId(999))

	gameId := model.NewGameId(42)
	require.Equal(t, model.GameQuestionId(420010000), model.NewGameQuestionId(gameId, 1))
	require.Equal(t, model.GameQuestionId(420020000), model.NewGameQuestionId(gameId, 2))
	require.Equal(t, model.GameQuestionId(420990000), model.NewGameQuestionId(gameId, 99))

	questionId := model.NewGameQuestionId(gameId, 17)
	require.Equal(t, model.GameAnswerId(420170100), model.NewGameAnswerId(questionId, 1))
	require.Equal(t, model.GameAnswerId(420170200), model.NewGameAnswerId(questionId, 2))
	require.Equal(t, model.GameAnswerId(420170900), model.NewGameAnswerId(questionId, 9))

	answerId := model.NewGameAnswerId(questionId, 3)
	require.Equal(t, model.GamePlayerAnswerId(420170301), model.NewGamePlayerAnswerId(answerId, model.NewGamePlayerId(1)))
	require.Equal(t, model.GamePlayerAnswerId(420170302), model.NewGamePlayerAnswerId(answerId, model.NewGamePlayerId(2)))
	require.Equal(t, model.GamePlayerAnswerId(420170309), model.NewGamePlayerAnswerId(answerId, model.NewGamePlayerId(9)))

	playerAnswerId := model.NewGamePlayerAnswerId(answerId, model.NewGamePlayerId(2))
	gameId, questionId, answerId, playerId := playerAnswerId.Split()
	require.Equal(t, model.GameId(420000000), gameId)
	require.Equal(t, model.GameQuestionId(420170000), questionId)
	require.Equal(t, model.GameAnswerId(420170300), answerId)
	require.Equal(t, model.GamePlayerId(2), playerId)

	gameId, questionId = answerId.Split()
	require.Equal(t, model.GameId(420000000), gameId)
	require.Equal(t, model.GameQuestionId(420170000), questionId)

	gameId = questionId.Split()
	require.Equal(t, model.GameId(420000000), gameId)
}
