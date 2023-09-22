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
// game service

type GameService interface {
	CreateGame(ctx context.Context, settings model.GameSettings) (*model.Game, error)
	RetrieveGame(ctx context.Context, id model.GameId) (*model.Game, error)
	DeleteGame(ctx context.Context, id model.GameId) error
}

func NewGameService(logger *zap.Logger, db *sql.DB, gameStore store.GameStore, gameQuestionStore store.GameQuestionStore, musicStore store.MusicStore, musiArtistStore store.MusicArtistStore, musicAlbumStore store.MusicAlbumStore, themeStore store.ThemeStore, themeQuestionStore store.ThemeQuestionStore) GameService {
	return &gameService{
		logger:             logger,
		db:                 db,
		gameStore:          gameStore,
		gameQuestionStore:  gameQuestionStore,
		musicStore:         musicStore,
		musiArtistStore:    musiArtistStore,
		musicAlbumStore:    musicAlbumStore,
		themeStore:         themeStore,
		themeQuestionStore: themeQuestionStore,
	}
}

type gameService struct {
	logger             *zap.Logger
	db                 *sql.DB
	gameStore          store.GameStore
	gameQuestionStore  store.GameQuestionStore
	musicStore         store.MusicStore
	musiArtistStore    store.MusicArtistStore
	musicAlbumStore    store.MusicAlbumStore
	themeStore         store.ThemeStore
	themeQuestionStore store.ThemeQuestionStore
}

func (s *gameService) CreateGame(ctx context.Context, settings model.GameSettings) (*model.Game, error) {

	var game *model.Game
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {

		var questions []*model.GameQuestion
		if settings.UseStore() {
			questions = s.createStoreQuestions(ctx, tx, settings)
		} else {
			questions = s.createLegacyQuestions(ctx, tx, settings)
		}

		game = &model.Game{
			Settings:  &settings,
			Players:   s.createPlayers(settings.NbPlayer),
			Questions: questions,
		}

		game = s.gameStore.Create(ctx, tx, game)

		for questionIndex, question := range game.Questions {
			question.Id = model.NewGameQuestionId(game.Id, questionIndex+1)
			for answerIndex, answer := range question.Answers {
				answer.Id = model.NewGameAnswerId(question.Id, answerIndex+1)
			}
		}
	})

	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *gameService) createLegacyQuestions(ctx context.Context, tx *sql.Tx, settings model.GameSettings) []*model.GameQuestion {
	return s.gameQuestionStore.SelectRandomQuestions(ctx, tx, settings)
}

func (s *gameService) createStoreQuestions(ctx context.Context, tx *sql.Tx, settings model.GameSettings) []*model.GameQuestion {

	//
	// select questions
	//

	s.logger.Info(fmt.Sprintf("[DEBUG] select %d questions", settings.NbQuestion))
	questions := s.themeQuestionStore.List(ctx, tx, &model.ThemeQuestionFilter{Random: true, Limit: settings.NbQuestion})

	//
	// retrieve themes
	//

	musics := map[model.MusicId]*model.Music{}
	artits := map[model.MusicArtistId]*model.MusicArtist{}
	albums := map[model.MusicAlbumId]*model.MusicAlbum{}
	themes := map[model.ThemeId]*model.Theme{}
	result := []*model.GameQuestion{}
	for _, question := range questions {

		//
		// retrieve theme
		//

		theme, found := themes[question.ThemeId]
		if !found {
			s.logger.Info(fmt.Sprintf("[DEBUG] retrieve theme %d", question.ThemeId))
			theme = s.themeStore.Retrieve(ctx, tx, question.ThemeId)
			s.logger.Info(fmt.Sprintf("[DEBUG] retrieve questions for theme %d", theme.Id))
			theme.Questions = s.themeQuestionStore.List(ctx, tx, &model.ThemeQuestionFilter{ThemeId: theme.Id})
		}

		//
		// retrieve music
		//

		music, found := musics[question.MusicId]
		if !found {
			s.logger.Info(fmt.Sprintf("[DEBUG] retrieve music %d", question.MusicId))
			music = s.musicStore.Retrieve(ctx, tx, question.MusicId)
			if music.ArtistId != 0 {
				music.Artist, found = artits[music.ArtistId]
				if !found {
					s.logger.Info(fmt.Sprintf("[DEBUG] retrieve artist %d", music.ArtistId))
					music.Artist = s.musiArtistStore.Retrieve(ctx, tx, music.ArtistId)
				}
			}
			if music.AlbumId != 0 {
				music.Album, found = albums[music.AlbumId]
				if !found {
					s.logger.Info(fmt.Sprintf("[DEBUG] retrieve album %d", music.AlbumId))
					music.Album = s.musicAlbumStore.Retrieve(ctx, tx, music.AlbumId)
				}
			}
		}

		//
		// select other answers
		//

		result = append(result, s.toQuestion(ctx, theme, question, music, settings.NbAnswer))
	}
	return result
}

func (s *gameService) toQuestion(ctx context.Context, theme *model.Theme, question *model.ThemeQuestion, music *model.Music, nbAnswer int) *model.GameQuestion {
	return &model.GameQuestion{
		Theme:   s.toTheme(ctx, theme),
		Music:   s.toMusic(ctx, music),
		Answers: s.toAnswers(ctx, theme, question, nbAnswer),
	}
}

func (s *gameService) toTheme(ctx context.Context, theme *model.Theme) *model.GameTheme {
	return &model.GameTheme{
		Title:  theme.Title,
		ImgUrl: theme.ImgUrl,
	}
}

func (s *gameService) toMusic(ctx context.Context, music *model.Music) *model.Music {
	return &model.Music{
		Id:     model.MusicId(music.Id),
		Name:   music.Name,
		Mp3Url: music.Mp3Url,
		Artist: s.toArtist(ctx, music.Artist),
		Album:  s.toAlbum(ctx, music.Album),
	}
}

func (s *gameService) toArtist(ctx context.Context, artist *model.MusicArtist) *model.MusicArtist {
	if artist == nil {
		return nil
	}
	return &model.MusicArtist{
		Id:       artist.Id,
		DeezerId: artist.DeezerId,
		Name:     artist.Name,
		ImgUrl:   artist.ImgUrl,
	}
}

func (s *gameService) toAlbum(ctx context.Context, album *model.MusicAlbum) *model.MusicAlbum {
	if album == nil {
		return nil
	}
	return &model.MusicAlbum{
		Id:       album.Id,
		DeezerId: album.DeezerId,
		Name:     album.Name,
		ImgUrl:   album.ImgUrl,
	}
}

func (s *gameService) toAnswers(ctx context.Context, theme *model.Theme, question *model.ThemeQuestion, nbAnswer int) []*model.GameAnswer {

	others := util.Filter(theme.Questions, func(other *model.ThemeQuestion) bool { return other.Id != question.Id })

	util.Shuffle(others)

	if len(others) > nbAnswer-1 {
		others = others[:nbAnswer-1]
	}

	answers := util.Convert(others, func(other *model.ThemeQuestion) *model.GameAnswer { return s.toAnswer(ctx, other, false) })
	answers = append(answers, s.toAnswer(ctx, question, true))

	util.Shuffle(answers)

	return answers
}

func (s *gameService) toAnswer(ctx context.Context, question *model.ThemeQuestion, correct bool) *model.GameAnswer {
	return &model.GameAnswer{
		Text:    question.Text,
		Hint:    question.Hint,
		Correct: correct,
	}
}

func (s *gameService) createPlayers(nbPlayer int) []*model.GamePlayer {
	players := make([]*model.GamePlayer, 0, nbPlayer)
	for playerNumber := 1; playerNumber <= nbPlayer; playerNumber++ {
		players = append(players, s.createPlayer(playerNumber))
	}
	return players
}

func (s *gameService) createPlayer(playerNumber int) *model.GamePlayer {
	return &model.GamePlayer{
		Id:     model.NewGamePlayerId(playerNumber),
		Name:   fmt.Sprintf("Player %02d", playerNumber),
		Active: true,
		Score:  0,
	}
}

func (s *gameService) RetrieveGame(ctx context.Context, id model.GameId) (*model.Game, error) {

	var game *model.Game
	err := util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		game = s.gameStore.Retrieve(ctx, tx, id)
	})
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *gameService) DeleteGame(ctx context.Context, id model.GameId) error {
	return util.SqlTransaction(ctx, s.db, func(tx *sql.Tx) {
		s.gameStore.Delete(ctx, tx, id)
	})
}
