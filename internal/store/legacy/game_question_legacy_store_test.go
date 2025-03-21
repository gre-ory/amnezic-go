package legacy_test

import (
	"context"
	"testing"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store/legacy"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLegacyMusicStore(t *testing.T) {
	ctx := context.Background()

	rootPath := "http://root"

	settings := model.GameSettings{
		Seed:       42,
		NbQuestion: 1,
		NbAnswer:   2,
		NbPlayer:   3,
		Sources: []model.Source{
			model.Source_Genre,
		},
	}
	logger := zap.L()

	store := legacy.NewGameQuestionLegacyStore(logger, rootPath)
	gotQuestions := store.SelectRandomQuestions(ctx, nil, settings)

	require.Equal(t, []*model.GameQuestion{
		{
			Id: 0,
			Theme: &model.GameTheme{
				Title: "Rock",
			},
			Music: &model.Music{
				Id:     3006021,
				Name:   "Time After Time",
				Mp3Url: model.Url("http://root/95532_EvaCassidy_TimeAfterTime_2000.mp3"),
				Artist: &model.MusicArtist{
					Name: "Eva Cassidy",
				},
			},
			Answers: []*model.GameAnswer{
				{
					Text:    "Eva Cassidy",
					Hint:    "Time After Time",
					Correct: true,
				},
				{
					Text: "Radiohead",
					Hint: "Kid A",
				},
			},
		},
	}, gotQuestions)
}

func Off_TestReproduce(t *testing.T) {
	ctx := context.Background()

	rootPath := "http://root"

	settings := model.GameSettings{
		Seed:       1678290075142,
		NbQuestion: 5,
		NbAnswer:   3,
		NbPlayer:   4,
		Sources:    []model.Source{model.Source_Legacy},
	}
	config := zap.NewDevelopmentConfig()
	config.Development = false
	logger, _ := config.Build()

	store := legacy.NewGameQuestionLegacyStore(logger, rootPath)
	gotQuestions := store.SelectRandomQuestions(ctx, nil, settings)

	require.Equal(t, &model.GameQuestion{
		Id: 2017,
		Theme: &model.GameTheme{
			Title: "Pop",
		},
		Music: &model.Music{
			Id:     2017,
			Name:   "Purple rain",
			Mp3Url: model.Url("http://root/6496_Prince_Purplerain_1984.mp3"),
		},
		Answers: []*model.GameAnswer{
			{
				Text: "Green Day",
				Hint: "Holiday",
			},
			{
				Text:    "Prince",
				Hint:    "Purple rain",
				Correct: true,
			},
			{
				Text: "Garbage",
				Hint: "I think I´m paranoid",
			},
		},
	}, gotQuestions[0])
}
