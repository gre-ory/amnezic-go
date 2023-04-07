package service

import (
	"context"

	"github.com/gre-ory/amnezic-go/internal/client"
	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// theme service

type ThemeService interface {
	Create(ctx context.Context, theme *model.Theme) (*model.Theme, error)
	Retrieve(ctx context.Context, id model.ThemeId) (*model.Theme, error)
	Update(ctx context.Context, theme *model.Theme) (*model.Theme, error)
	Delete(ctx context.Context, id model.ThemeId) error
}

func NewThemeService(logger *zap.Logger, deezerClient client.DeezerClient, themeStore store.ThemeStore, themequestionStore store.ThemeQuestionStore, musicStore         store.MusicStore) ThemeService {
	return &themeService{
		logger:             logger,
		deezerClient:       deezerClient,
		themeStore:         themeStore,
		themeQuestionStore: themequestionStore,
		musicStore: musicStore,
	}
}

type themeService struct {
	logger             *zap.Logger
	deezerClient       client.DeezerClient
	themeStore         store.ThemeStore
	themeQuestionStore store.ThemeQuestionStore
	musicStore         store.MusicStore
}

// //////////////////////////////////////////////////
// create

func (s *themeService) Create(ctx context.Context, theme *model.Theme) (*model.Theme, error) {

	//
	// create theme
	//

	created, err := s.themeStore.Create(ctx, theme.Copy())
	if err != nil {
		return nil, err
	}

	//
	// create questions
	//

	for _, question := range theme.Questions {
		createdQuestion, err := s.themeQuestionStore.Create(ctx, question)
		if err != nil {
			return nil, err
		}
		created.Questions = append(created.Questions, createdQuestion)
	}

	return created, nil
}

// //////////////////////////////////////////////////
// retrieve

func (s *themeService) Retrieve(ctx context.Context, id model.ThemeId) (*model.Theme, error) {

	//
	// retreve theme
	//

	theme, err := s.themeStore.Retrieve(ctx, id)
	if err != nil {
		return nil, err
	}

	//
	// retrieve theme questions
	//

	theme.Questions, err = s.themeQuestionStore.List(ctx, &model.ThemeQuestionFilter{ThemeId: id})
	if err != nil {
		return nil, err
	}

	return theme, nil
}

// //////////////////////////////////////////////////
// update

func (s *themeService) Update(ctx context.Context, theme *model.Theme) (*model.Theme, error) {

	//
	// update theme
	//

	orig, err := s.themeStore.Retrieve(ctx, theme.Id)
	if err != nil {
		return nil, err
	}

	updated := theme.Copy()
	if !orig.Equal(updated) {
		updated, err = s.themeStore.Update(ctx, updated)
		if err != nil {
			return nil, err
		}
	}

	//
	// update questions
	//

	origQuestions, err := s.themeQuestionStore.List(ctx, &model.ThemeQuestionFilter{ThemeId: theme.Id})
	if err != nil {
		return nil, err
	}

	for _, origQuestion := range origQuestions {
		_, found := util.FindIf(theme.Questions, func(question *model.ThemeQuestion) bool { return question.Id == origQuestion.Id })
		if !found {
			s.themeQuestionStore.Delete(ctx, &model.ThemeQuestionFilter{ThemeQuestionId: origQuestion.Id})
		}
	}
	for _, question := range theme.Questions {
		_, found := util.FindIf(origQuestions, func(origQuestion *model.ThemeQuestion) bool { return question.Id == origQuestion.Id })
		if !found {
			createdQuestion, err := s.themeQuestionStore.Create(ctx, question)
			if err != nil {
				return nil, err
			}
			updated.Questions = append(updated.Questions, createdQuestion)
		} else {
			updatedQuestion, err := s.themeQuestionStore.Update(ctx, question)
			if err != nil {
				return nil, err
			}
			updated.Questions = append(updated.Questions, updatedQuestion)
		}
	}

	return updated, nil
}

// //////////////////////////////////////////////////
// delete

func (s *themeService) Delete(ctx context.Context, id model.ThemeId) error {

	//
	// delete theme
	//

	err := s.themeStore.Delete(ctx, &model.ThemeFilter{ThemeId: id})
	if err != nil {
		return err
	}

	//
	// delete questions
	//

	err = s.themeQuestionStore.Delete(ctx, &model.ThemeQuestionFilter{ThemeId: id})
	if err != nil {
		return err
	}

	return nil
}
