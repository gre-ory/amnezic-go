package store_test

import (
	"context"
	"testing"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/stretchr/testify/require"
)

func TestLegacyMusicStore(t *testing.T) {
	ctx := context.Background()

	rootPath := "http://root"

	settings := model.GameSettings{
		Seed:       42,
		UseLegacy:  true,
		NbQuestion: 2,
		NbAnswer:   3,
		NbPlayer:   4,
	}

	store := store.NewLegacyMusicStore(rootPath)
	gotQuestions, gotErr := store.SelectRandomQuestions(ctx, settings)

	require.NoError(t, gotErr)
	require.Equal(t, []*model.Question{
		{
			Id: 2017,
			Theme: model.Theme{
				Title: "Pop",
			},
			Music: model.Music{
				Id:     2017,
				Name:   "Purple rain",
				Mp3Url: "http://root/6496_Prince_Purplerain_1984.mp3",
			},
			Answers: []*model.Answer{
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
			Theme: model.Theme{
				Title: "Nostalgie",
			},
			Music: model.Music{
				Id:     5001,
				Name:   "Dejeuner en paix",
				Mp3Url: "http://root/Rem2-08-01.mp3",
			},
			Answers: []*model.Answer{
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
