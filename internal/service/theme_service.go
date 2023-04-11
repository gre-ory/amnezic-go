package service

import (
	"context"
	"fmt"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// theme service

type ThemeService interface {
	CreateTheme(ctx context.Context, theme *model.Theme) (*model.Theme, error)
	RetrieveTheme(ctx context.Context, id model.ThemeId) (*model.Theme, error)
	UpdateTheme(ctx context.Context, theme *model.Theme) (*model.Theme, error)
	DeleteTheme(ctx context.Context, id model.ThemeId) error
}

func NewThemeService(logger *zap.Logger, themeStore store.ThemeStore, themequestionStore store.ThemeQuestionStore, musicStore store.MusicStore) ThemeService {
	return &themeService{
		logger:             logger,
		themeStore:         themeStore,
		themeQuestionStore: themequestionStore,
		musicStore:         musicStore,
	}
}

type themeService struct {
	logger             *zap.Logger
	themeStore         store.ThemeStore
	themeQuestionStore store.ThemeQuestionStore
	musicStore         store.MusicStore
}

// //////////////////////////////////////////////////
// create

func (s *themeService) CreateTheme(ctx context.Context, theme *model.Theme) (*model.Theme, error) {

	var created *model.Theme
	var err error

	defer func() {
		if err == nil {
			s.logger.Info(fmt.Sprintf("[ OK ] create theme: %#v", theme))
		} else {
			s.logger.Info(fmt.Sprintf("[ KO ] create theme: %#v", theme), zap.Error(err))
		}
	}()

	//
	// create theme
	//

	s.logger.Info(fmt.Sprintf("[DEBUG] create theme: %#v", theme.Copy()))
	created, err = s.themeStore.Create(ctx, theme)
	if err != nil {
		return nil, err
	}

	//
	// create questions
	//

	// created.Questions = make([]*model.ThemeQuestion, 0, len(theme.Questions))
	for _, question := range theme.Questions {
		question.ThemeId = created.Id
		s.logger.Info(fmt.Sprintf("[DEBUG] create question: %#v", question.Copy()))
		createdQuestion, err := s.themeQuestionStore.Create(ctx, question)
		if err != nil {
			return nil, err
		}
		created.Questions = append(created.Questions, createdQuestion)
	}

	//
	// attach musics
	//

	created.Questions = util.Convert(created.Questions, s.AttachMusic(ctx))

	return created, nil
}

// //////////////////////////////////////////////////
// retrieve

func (s *themeService) RetrieveTheme(ctx context.Context, id model.ThemeId) (*model.Theme, error) {

	var theme *model.Theme
	var err error

	defer func() {
		if err == nil {
			s.logger.Info(fmt.Sprintf("[ OK ] retrieve theme %d", id))
		} else {
			s.logger.Info(fmt.Sprintf("[ KO ] retrieve theme %d", id), zap.Error(err))
		}
	}()

	//
	// retreve theme
	//

	s.logger.Info(fmt.Sprintf("[DEBUG] retrieve theme %d", id))
	theme, err = s.themeStore.Retrieve(ctx, id)
	if err != nil {
		return nil, err
	}

	//
	// retrieve theme questions
	//

	s.logger.Info(fmt.Sprintf("[DEBUG] retrieve theme %d", theme.Id))
	theme.Questions, err = s.themeQuestionStore.List(ctx, &model.ThemeQuestionFilter{ThemeId: id})
	if err != nil {
		return nil, err
	}

	//
	// attach musics
	//

	theme.Questions = util.Convert(theme.Questions, s.AttachMusic(ctx))

	return theme, nil
}

// //////////////////////////////////////////////////
// update

func (s *themeService) UpdateTheme(ctx context.Context, theme *model.Theme) (*model.Theme, error) {

	var updated *model.Theme
	var err error

	defer func() {
		if err == nil {
			s.logger.Info(fmt.Sprintf("[ OK ] update theme: %#v", theme))
		} else {
			s.logger.Info(fmt.Sprintf("[ KO ] update theme: %#v", theme), zap.Error(err))
		}
	}()

	//
	// retrieve theme
	//

	orig, err := s.RetrieveTheme(ctx, theme.Id)
	if err != nil {
		return nil, err
	}

	//
	// update theme
	//

	s.logger.Info(fmt.Sprintf("[DEBUG] update theme: %#v", theme.Copy()))
	updated, err = s.themeStore.Update(ctx, theme)
	if err != nil {
		return nil, err
	}

	//
	// delete questions
	//

	for _, origQuestion := range orig.Questions {
		_, found := util.FindIf(theme.Questions, func(question *model.ThemeQuestion) bool { return question.Id != 0 && question.Id == origQuestion.Id })
		if !found {
			s.logger.Info(fmt.Sprintf("[DEBUG] delete question: %#v", origQuestion.Copy()))
			err := s.themeQuestionStore.Delete(ctx, &model.ThemeQuestionFilter{ThemeQuestionId: origQuestion.Id})
			if err != nil {
				return nil, err
			}
		}
	}

	//
	// upsert questions
	//

	for _, question := range theme.Questions {
		question.ThemeId = updated.Id
		_, found := util.FindIf(orig.Questions, func(origQuestion *model.ThemeQuestion) bool {
			return question.Id != 0 && question.Id == origQuestion.Id
		})
		if found {
			s.logger.Info(fmt.Sprintf("[DEBUG] update question: %#v", question.Copy()))
			updatedQuestion, err := s.themeQuestionStore.Update(ctx, question)
			if err != nil {
				return nil, err
			}
			updated.Questions = append(updated.Questions, updatedQuestion)
		} else {
			s.logger.Info(fmt.Sprintf("[DEBUG] create question: %#v", question.Copy()))
			createdQuestion, err := s.themeQuestionStore.Create(ctx, question)
			if err != nil {
				return nil, err
			}
			updated.Questions = append(updated.Questions, createdQuestion)
		}
	}

	//
	// attach musics
	//

	updated.Questions = util.Convert(updated.Questions, s.AttachMusic(ctx))

	return updated, nil
}

// //////////////////////////////////////////////////
// delete

func (s *themeService) DeleteTheme(ctx context.Context, id model.ThemeId) error {

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

// //////////////////////////////////////////////////
// attach music

func (s *themeService) AttachMusic(ctx context.Context) func(question *model.ThemeQuestion) *model.ThemeQuestion {
	return func(question *model.ThemeQuestion) *model.ThemeQuestion {
		if question.MusicId != 0 {
			question.Music, _ = s.musicStore.Retrieve(ctx, question.MusicId)
		}
		return question
	}
}
