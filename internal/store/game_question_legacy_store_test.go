package store_test

import (
	"context"
	"testing"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLegacyMusicStore(t *testing.T) {
	ctx := context.Background()

	rootPath := "http://root"

	settings := model.GameSettings{
		Seed:       42,
		NbQuestion: 2,
		NbAnswer:   3,
		NbPlayer:   4,
	}
	logger := zap.L()

	store := store.NewLegacyMusicStore(logger, rootPath)
	gotQuestions, gotErr := store.SelectRandomQuestions(ctx, settings)

	require.NoError(t, gotErr)
	require.Equal(t, []*model.GameQuestion{
		{
			Id: 2017,
			Theme: &model.GameTheme{
				Title: "Pop",
			},
			Music: &model.Music{
				Id:     2017,
				Name:   "Purple rain",
				Mp3Url: "http://root/6496_Prince_Purplerain_1984.mp3",
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
		},
		{
			Id: 5001,
			Theme: &model.GameTheme{
				Title: "Nostalgie",
			},
			Music: &model.Music{
				Id:     5001,
				Name:   "Dejeuner en paix",
				Mp3Url: "http://root/Rem2-08-01.mp3",
			},
			Answers: []*model.GameAnswer{
				{
					Text: "Claude francois",
					Hint: "Cette annee la",
				},
				{
					Text: "Florent Pagny",
					Hint: "Si tu veux m'essayer",
				},
				{
					Text:    "Stephan Eicher",
					Hint:    "Dejeuner en paix",
					Correct: true,
				},
			},
		},
	}, gotQuestions)
}

func TestReproduce(t *testing.T) {
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

	store := store.NewLegacyMusicStore(logger, rootPath)
	gotQuestions, gotErr := store.SelectRandomQuestions(ctx, settings)

	require.NoError(t, gotErr)
	require.Equal(t, &model.GameQuestion{
		Id: 2017,
		Theme: &model.GameTheme{
			Title: "Pop",
		},
		Music: &model.Music{
			Id:     2017,
			Name:   "Purple rain",
			Mp3Url: "http://root/6496_Prince_Purplerain_1984.mp3",
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
