package legacy

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/store"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"

	_ "embed"
)

// //////////////////////////////////////////////////
// legacy music store

const (
	RootPath_FreeDotFr = "http://gregory.valigiani.free.fr/amnezic/"
)

func NewGameQuestionLegacyStore(logger *zap.Logger, rootPath string) store.GameQuestionStore {
	store := &gameQuestionLegacyMusicStore{
		logger:           logger,
		rootPath:         strings.TrimRight(rootPath, "/"),
		mediaIdsBySource: make(map[model.Source][]int64, 10),
		media:            make(map[int64]*JsonLegacyMedia, 2000),
		genres:           make(map[int64]*JsonLegacyGenre, 200),
	}
	store.Load()
	return store
}

// //////////////////////////////////////////////////
// implementation

type gameQuestionLegacyMusicStore struct {
	logger           *zap.Logger
	rootPath         string
	mediaIdsBySource map[model.Source][]int64
	media            map[int64]*JsonLegacyMedia
	genres           map[int64]*JsonLegacyGenre
}

func (s *gameQuestionLegacyMusicStore) SelectRandomQuestions(ctx context.Context, _ *sql.Tx, settings model.GameSettings) []*model.GameQuestion {

	//
	// validate
	//

	if settings.NbQuestion <= 0 {
		panic(model.ErrInvalidNumberOfQuestion)
	}
	if settings.NbAnswer <= 0 {
		panic(model.ErrInvalidNumberOfAnswer)
	}

	//
	// random seed
	//

	rand.Seed(settings.Seed)

	//
	// select & shuffle media ids
	//

	mediaIds := make([]int64, 0, 2000)
	for _, source := range settings.Sources {
		mediaIds = append(mediaIds, s.mediaIdsBySource[source]...)
	}
	util.Shuffle(mediaIds)

	//
	// select subset
	//

	if len(mediaIds) > settings.NbQuestion {
		mediaIds = mediaIds[:settings.NbQuestion]
	}

	//
	// building questions
	//

	questions := make([]*model.GameQuestion, 0, settings.NbQuestion)
	for _, mediaId := range mediaIds {
		media := s.media[mediaId]
		genre := s.genres[media.GenreId]
		questions = append(questions, s.toQuestion(ctx, genre, media, settings.NbAnswer))
	}

	return questions
}

func (s *gameQuestionLegacyMusicStore) toQuestion(ctx context.Context, genre *JsonLegacyGenre, media *JsonLegacyMedia, nbAnswer int) *model.GameQuestion {
	return &model.GameQuestion{
		Theme:   s.toTheme(ctx, genre),
		Music:   s.toMusic(ctx, media),
		Answers: s.toAnswers(ctx, genre, media, nbAnswer),
	}
}

func (s *gameQuestionLegacyMusicStore) toTheme(ctx context.Context, genre *JsonLegacyGenre) *model.GameTheme {
	return &model.GameTheme{
		Title: genre.Genre,
	}
}

func (s *gameQuestionLegacyMusicStore) toMusic(ctx context.Context, media *JsonLegacyMedia) *model.Music {
	return &model.Music{
		Id:     model.MusicId(media.Id),
		Name:   media.Title,
		Mp3Url: s.toMp3Url(ctx, media),
		Artist: s.toArtist(ctx, media.Artist),
	}
}

func (s *gameQuestionLegacyMusicStore) toMp3Url(ctx context.Context, media *JsonLegacyMedia) model.Url {
	if s.rootPath != "" {
		return model.Url(fmt.Sprintf("%s/%s", s.rootPath, media.MusicFileName))
	}
	return model.Url(media.MusicFileName)
}

func (s *gameQuestionLegacyMusicStore) toArtist(ctx context.Context, artist *JsonLegacyArtist) *model.MusicArtist {
	if artist == nil {
		return nil
	}
	return &model.MusicArtist{
		Name: artist.Name,
	}
}

func (s *gameQuestionLegacyMusicStore) toAnswers(ctx context.Context, genre *JsonLegacyGenre, media *JsonLegacyMedia, nbAnswer int) []*model.GameAnswer {

	others := util.Filter(genre.Media, func(other *JsonLegacyMedia) bool { return other.Id != media.Id })

	util.Shuffle(others)

	if len(others) > nbAnswer-1 {
		others = others[:nbAnswer-1]
	}

	answers := util.Convert(others, func(other *JsonLegacyMedia) *model.GameAnswer { return s.toAnswer(ctx, other, false) })
	answers = append(answers, s.toAnswer(ctx, media, true))

	util.Shuffle(answers)

	return answers
}

func (s *gameQuestionLegacyMusicStore) toAnswer(ctx context.Context, media *JsonLegacyMedia, correct bool) *model.GameAnswer {
	if media.Artist == nil {
		return &model.GameAnswer{
			Text:    media.Title,
			Correct: correct,
		}
	}
	return &model.GameAnswer{
		Text:    media.Artist.Name,
		Hint:    media.Title,
		Correct: correct,
	}
}

// //////////////////////////////////////////////////
// load

//go:embed resources/decade.json
var decadeJsonBytes []byte

//go:embed resources/genre.json
var genreJsonBytes []byte

//go:embed resources/legacy.json
var legacyJsonBytes []byte

func (s *gameQuestionLegacyMusicStore) Load() *gameQuestionLegacyMusicStore {

	s.LoadSource(model.Source_Legacy, legacyJsonBytes)
	s.LoadSource(model.Source_Decade, decadeJsonBytes)
	s.LoadSource(model.Source_Genre, genreJsonBytes)

	return s
}

func (s *gameQuestionLegacyMusicStore) LoadSource(source model.Source, embedBytes []byte) *gameQuestionLegacyMusicStore {

	jsonLegacy := JsonLegacy{}
	if err := json.Unmarshal(embedBytes, &jsonLegacy); err != nil {
		panic(err)
	}

	sourceIndex := len(s.mediaIdsBySource)
	sourceId := int64(1000 * 1000 * (sourceIndex + 1))

	mediaIds := make([]int64, 0, 1000)

	for genreIndex, genre := range jsonLegacy.Genres {
		genre.Id = sourceId + int64(1000*(genreIndex+1))
		s.genres[genre.Id] = genre
		for mediaIndex, media := range genre.Media {
			media.GenreId = genre.Id
			media.Id = genre.Id + int64(mediaIndex) + 1
			mediaIds = append(mediaIds, media.Id)
			s.media[media.Id] = media
		}
	}

	s.mediaIdsBySource[source] = mediaIds

	return s
}

type JsonLegacy struct {
	Genres []*JsonLegacyGenre `json:"genres,omitempty"`
}

type JsonLegacyGenre struct {
	Id    int64              `json:"id,omitempty"`
	Genre string             `json:"genre,omitempty"`
	Media []*JsonLegacyMedia `json:"media,omitempty"`
}

type JsonLegacyMedia struct {
	Id            int64             `json:"id,omitempty"`
	GenreId       int64             `json:"genreId,omitempty"`
	Title         string            `json:"title,omitempty"`
	MusicFileName string            `json:"music,omitempty"`
	Artist        *JsonLegacyArtist `json:"artist,omitempty"`
}

type JsonLegacyArtist struct {
	Name string `json:"name,omitempty"`
}
