package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// theme service

type ThemeService interface {
	ListThemes(ctx context.Context) ([]*model.ThemeInfo, error)
	CreateTheme(ctx context.Context, theme *model.Theme) (*model.Theme, error)
	RetrieveTheme(ctx context.Context, id model.ThemeId) (*model.Theme, error)
	UpdateTheme(ctx context.Context, theme *model.Theme) (*model.Theme, error)
	DeleteTheme(ctx context.Context, id model.ThemeId) error
	AddQuestion(ctx context.Context, question *model.ThemeQuestion) (*model.Theme, error)
	UpdateQuestion(ctx context.Context, question *model.ThemeQuestion) (*model.Theme, error)
	RemoveQuestion(ctx context.Context, id model.ThemeId, questionId model.ThemeQuestionId) (*model.Theme, error)
}

func NewThemeService(logger *zap.Logger, db *sql.DB, themeStore store.ThemeStore, themequestionStore store.ThemeQuestionStore, musicStore store.MusicStore, artistStore store.MusicArtistStore, albumStore store.MusicAlbumStore) ThemeService {
	return &themeService{
		logger:             logger,
		db:                 db,
		themeStore:         themeStore,
		themeQuestionStore: themequestionStore,
		musicStore:         musicStore,
		artistStore:        artistStore,
		albumStore:         albumStore,
	}
}

type themeService struct {
	logger             *zap.Logger
	db                 *sql.DB
	themeStore         store.ThemeStore
	themeQuestionStore store.ThemeQuestionStore
	musicStore         store.MusicStore
	artistStore        store.MusicArtistStore
	albumStore         store.MusicAlbumStore
}

func (s *themeService) ListThemes(ctx context.Context) ([]*model.ThemeInfo, error) {

	var themes []*model.Theme
	var infos []*model.ThemeInfo
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		//
		// list themes
		//

		s.logger.Info("[DEBUG] list themes")
		themes = s.themeStore.List(ctx, tx, nil)
		infos = util.Convert(themes, (*model.Theme).GetInfo)

		//
		// count questions
		//

		s.logger.Info("[DEBUG] count questions")
		count := s.themeQuestionStore.CountByTheme(ctx, tx)
		for _, info := range infos {
			info.NbQuestion = count[info.Id]
		}
	})

	if err != nil {
		s.logger.Info("[ KO ] list themes", zap.Error(err))
		return nil, err
	}
	s.logger.Info("[ OK ] list themes")
	return infos, nil
}

// //////////////////////////////////////////////////
// create

func (s *themeService) CreateTheme(ctx context.Context, theme *model.Theme) (*model.Theme, error) {

	var created *model.Theme
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// create theme
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] create theme: %#v", theme.Copy()))
		created = s.themeStore.Create(ctx, tx, theme)

		//
		// create questions
		//

		// created.Questions = make([]*model.ThemeQuestion, 0, len(theme.Questions))
		for _, question := range theme.Questions {
			question.ThemeId = created.Id
			s.logger.Info(fmt.Sprintf("[DEBUG] create question: %#v", question.Copy()))
			createdQuestion := s.themeQuestionStore.Create(ctx, tx, question)
			created.Questions = append(created.Questions, createdQuestion)
		}

		//
		// attach musics
		//

		created.Questions = util.Convert(created.Questions, s.AttachMusic(ctx, tx))

	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] create theme: %#v", theme), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] create theme: %#v", created))
	return created, nil
}

// //////////////////////////////////////////////////
// retrieve

func (s *themeService) RetrieveTheme(ctx context.Context, id model.ThemeId) (*model.Theme, error) {

	var theme *model.Theme
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		theme = s.retrieveTheme(ctx, tx, id)
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] retrieve theme %d", id), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] retrieve theme %d", id))
	return theme, nil
}

func (s *themeService) retrieveTheme(ctx context.Context, tx *sql.Tx, id model.ThemeId) *model.Theme {

	//
	// retreve theme
	//

	s.logger.Info(fmt.Sprintf("[DEBUG] retrieve theme %d", id))
	theme := s.themeStore.Retrieve(ctx, tx, id)

	//
	// retrieve theme questions
	//

	s.logger.Info(fmt.Sprintf("[DEBUG] retrieve questions for theme %d", theme.Id))
	filter := &model.ThemeQuestionFilter{
		ThemeIds: []model.ThemeId{id},
	}
	theme.Questions = s.themeQuestionStore.List(ctx, tx, filter)

	//
	// attach musics
	//

	theme.Questions = util.Convert(theme.Questions, s.AttachMusic(ctx, tx))

	return theme
}

// //////////////////////////////////////////////////
// update

func (s *themeService) UpdateTheme(ctx context.Context, theme *model.Theme) (*model.Theme, error) {

	var updated *model.Theme
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// retrieve theme
		//

		orig := s.retrieveTheme(ctx, tx, theme.Id)

		//
		// update theme
		//

		s.logger.Info(fmt.Sprintf("[DEBUG] update theme: %#v", theme.Copy()))
		updated = s.themeStore.Update(ctx, tx, theme)

		//
		// delete questions
		//

		for _, origQuestion := range orig.Questions {
			_, found := util.FindIf(theme.Questions, func(question *model.ThemeQuestion) bool { return question.Id != 0 && question.Id == origQuestion.Id })
			if !found {
				s.logger.Info(fmt.Sprintf("[DEBUG] delete question: %#v", origQuestion.Copy()))
				s.themeQuestionStore.Delete(ctx, tx, &model.ThemeQuestionFilter{ThemeQuestionId: origQuestion.Id})
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
				updatedQuestion := s.themeQuestionStore.Update(ctx, tx, question)
				updated.Questions = append(updated.Questions, updatedQuestion)
			} else {
				s.logger.Info(fmt.Sprintf("[DEBUG] create question: %#v", question.Copy()))
				createdQuestion := s.themeQuestionStore.Create(ctx, tx, question)
				updated.Questions = append(updated.Questions, createdQuestion)
			}
		}

		//
		// attach musics
		//

		updated.Questions = util.Convert(updated.Questions, s.AttachMusic(ctx, tx))
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] update theme %d - %s", theme.Id, theme.Title), zap.Object("theme", theme), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] update theme %d - %s", updated.Id, updated.Title), zap.Object("theme", updated))
	return updated, nil

}

// //////////////////////////////////////////////////
// delete

func (s *themeService) DeleteTheme(ctx context.Context, id model.ThemeId) error {

	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// delete theme
		//

		s.themeStore.Delete(ctx, tx, &model.ThemeFilter{ThemeId: id})

		//
		// delete questions
		//

		filter := &model.ThemeQuestionFilter{
			ThemeIds: []model.ThemeId{id},
		}
		s.themeQuestionStore.Delete(ctx, tx, filter)

	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] delete theme: %#v", id), zap.Error(err))
		return err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] delete theme: %#v", id))
	return nil
}

// //////////////////////////////////////////////////
// add question

func (s *themeService) AddQuestion(ctx context.Context, question *model.ThemeQuestion) (*model.Theme, error) {

	if question == nil {
		return nil, model.ErrInvalidThemeQuestion
	}
	if err := question.Validate(); err != nil {
		return nil, err
	}

	var created *model.ThemeQuestion
	var theme *model.Theme
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// check for duplicate
		//

		alreadyExists := s.themeQuestionStore.IsMusicInTheme(ctx, tx, question.ThemeId, question.MusicId)
		if alreadyExists {
			panic(model.ErrMusicAlreadyInTheme)
		}

		//
		// create question
		//

		created = s.themeQuestionStore.Create(ctx, tx, question)

		//
		// retrieve theme
		//

		theme = s.retrieveTheme(ctx, tx, created.ThemeId)
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] add question to theme %d", question.ThemeId), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] add question %d to theme %d", created.Id, created.ThemeId), zap.Object("theme", theme))
	return theme, nil
}

// //////////////////////////////////////////////////
// update question

func (s *themeService) UpdateQuestion(ctx context.Context, question *model.ThemeQuestion) (*model.Theme, error) {

	if question == nil {
		return nil, model.ErrInvalidThemeQuestion
	}
	if err := question.Validate(); err != nil {
		return nil, err
	}

	var updated *model.ThemeQuestion
	var theme *model.Theme
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// check for invariant
		//

		orig := s.themeQuestionStore.Retrieve(ctx, tx, question.Id)
		if orig.ThemeId != question.ThemeId {
			panic(model.ErrCouldNotUpdateThemeId)
		}
		if orig.MusicId != question.MusicId {
			panic(model.ErrCouldNotUpdateMusicId)
		}

		//
		// update question
		//

		updated = s.themeQuestionStore.Update(ctx, tx, question)

		//
		// retrueve theme
		//

		theme = s.retrieveTheme(ctx, tx, updated.ThemeId)
	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] update question in theme %d", question.ThemeId), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] update question %d in theme %d", updated.Id, updated.ThemeId), zap.Object("theme", theme))
	return theme, nil
}

// //////////////////////////////////////////////////
// remove question

func (s *themeService) RemoveQuestion(ctx context.Context, id model.ThemeId, questionId model.ThemeQuestionId) (*model.Theme, error) {

	var theme *model.Theme
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		//
		// delete question
		//

		s.themeQuestionStore.Delete(ctx, tx, &model.ThemeQuestionFilter{ThemeQuestionId: questionId})

		//
		// retrieve theme
		//

		theme = s.retrieveTheme(ctx, tx, id)

	})

	if err != nil {
		s.logger.Info(fmt.Sprintf("[ KO ] delete question %d from theme %d", questionId, id), zap.Error(err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("[ OK ] delete question %d from theme %d", questionId, id), zap.Object("theme", theme))
	return theme, nil
}

// //////////////////////////////////////////////////
// attach music

func (s *themeService) AttachMusic(ctx context.Context, tx *sql.Tx) func(question *model.ThemeQuestion) *model.ThemeQuestion {

	// TODO avoid multiple queries

	return func(question *model.ThemeQuestion) *model.ThemeQuestion {
		if question.MusicId != 0 {
			question.Music = s.musicStore.Retrieve(ctx, tx, question.MusicId)
			if question.Music != nil {
				if question.Music.AlbumId != 0 {
					question.Music.Album = s.albumStore.Retrieve(ctx, tx, question.Music.AlbumId)
				}
				if question.Music.ArtistId != 0 {
					question.Music.Artist = s.artistStore.Retrieve(ctx, tx, question.Music.ArtistId)
				}
			}
		}
		return question
	}
}
