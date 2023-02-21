package store

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/util"

	_ "embed"
)

// //////////////////////////////////////////////////
// legacy music store

const (
	RootPath_AmnezicDotCom = "http://cache.amnezic.com/"
	RootPath_FreeDotFr     = "http://gregory.valigiani.free.fr/amnezic/"
)

func NewLegacyMusicStore(rootPath string) MusicStore {
	store := &legacyMusicStore{
		rootPath: strings.TrimRight(rootPath, "/"),
		mediaIds: make([]int64, 0, 2000),
		media:    make(map[int64]JsonLegacyMedia, 2000),
		genres:   make(map[int64]JsonLegacyGenre, 200),
	}
	store.Load()
	return store
}

// //////////////////////////////////////////////////
// implementation

type legacyMusicStore struct {
	rootPath string
	mediaIds []int64
	media    map[int64]JsonLegacyMedia
	genres   map[int64]JsonLegacyGenre
}

func (s *legacyMusicStore) SelectRandomQuestions(ctx context.Context, settings model.GameSettings) ([]*model.Question, error) {

	//
	// validate
	//

	if settings.NbQuestion <= 0 {
		return nil, fmt.Errorf("invalid number of question")
	}
	if settings.NbAnswer <= 0 {
		return nil, fmt.Errorf("invalid number of answer")
	}

	//
	// random seed
	//

	rand.Seed(settings.Seed)

	//
	// copy & shuffle media ids
	//

	mediaIds := s.mediaIds
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

	questions := make([]*model.Question, 0, settings.NbQuestion)
	for _, mediaId := range mediaIds {
		media := s.media[mediaId]
		genre := s.genres[media.GenreId]
		questions = append(questions, s.toQuestion(ctx, genre, media, settings.NbAnswer))
	}

	return questions, nil
}

func (s *legacyMusicStore) toQuestion(ctx context.Context, genre JsonLegacyGenre, media JsonLegacyMedia, nbAnswer int) *model.Question {
	return &model.Question{
		Theme:   s.toTheme(ctx, genre),
		Music:   s.toMusic(ctx, media),
		Answers: s.toAnswers(ctx, genre, media, nbAnswer),
	}
}

func (s *legacyMusicStore) toTheme(ctx context.Context, genre JsonLegacyGenre) model.Theme {
	return model.Theme{
		Title: genre.Genre,
	}
}

func (s *legacyMusicStore) toMusic(ctx context.Context, media JsonLegacyMedia) model.Music {
	return model.Music{
		Id:     media.Id,
		Name:   media.Title,
		Mp3Url: s.toMp3Url(ctx, media),
	}
}

func (s *legacyMusicStore) toMp3Url(ctx context.Context, media JsonLegacyMedia) string {
	if s.rootPath != "" {
		return fmt.Sprintf("%s/%s", s.rootPath, media.MusicFileName)
	}
	return media.MusicFileName
}

func (s *legacyMusicStore) toAnswers(ctx context.Context, genre JsonLegacyGenre, media JsonLegacyMedia, nbAnswer int) []*model.Answer {

	others := util.Filter(genre.Media, func(other JsonLegacyMedia) bool { return other.Id != media.Id })
	util.Shuffle(others)
	if len(others) > nbAnswer-1 {
		others = others[:nbAnswer-1]
	}

	answers := util.Convert(others, func(other JsonLegacyMedia) *model.Answer { return s.toAnswer(ctx, other, false) })
	answers = append(answers, s.toAnswer(ctx, media, true))

	util.Shuffle(answers)

	return answers
}

func (s *legacyMusicStore) toAnswer(ctx context.Context, media JsonLegacyMedia, correct bool) *model.Answer {
	return &model.Answer{
		Text:    media.Artist.Name,
		Hint:    media.Title,
		Correct: correct,
	}
}

// //////////////////////////////////////////////////
// load

//go:embed resources/legacy.json
var legacyJsonBytes []byte

func (s *legacyMusicStore) Load() *legacyMusicStore {

	jsonLegacy := JsonLegacy{}
	if err := json.Unmarshal(legacyJsonBytes, &jsonLegacy); err != nil {
		panic(err)
	}

	for genreIndex, genre := range jsonLegacy.Genres {
		genre.Id = int64(1000 * (genreIndex + 1))
		s.genres[genre.Id] = genre
		for mediaIndex, media := range genre.Media {
			media.GenreId = genre.Id
			media.Id = genre.Id + int64(mediaIndex) + 1
			s.mediaIds = append(s.mediaIds, media.Id)
			s.media[media.Id] = media
		}
	}

	return s
}

type JsonLegacy struct {
	Genres []JsonLegacyGenre `json:"genres,omitempty"`
}

type JsonLegacyGenre struct {
	Id    int64             `json:"id,omitempty"`
	Genre string            `json:"genre,omitempty"`
	Media []JsonLegacyMedia `json:"media,omitempty"`
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
