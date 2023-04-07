package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/service"
	"github.com/gre-ory/amnezic-go/internal/util"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// game handler

func NewGamehandler(logger *zap.Logger, service service.GameService) Handler {
	return &gameHandler{
		logger:  logger,
		service: service,
	}
}

type gameHandler struct {
	logger  *zap.Logger
	service service.GameService
}

// //////////////////////////////////////////////////
// register

func (h *gameHandler) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPut, "/api/game/new", h.handleCreateGame)
	router.HandlerFunc(http.MethodGet, "/api/game/:game_id", h.handleRetrieveGame)
	router.HandlerFunc(http.MethodPost, "/api/game/:game_id", h.handleUpdateGame)
	router.HandlerFunc(http.MethodDelete, "/api/game/:game_id", h.handleDeleteGame)
}

// //////////////////////////////////////////////////
// decode

func (h *gameHandler) extractStringParameter(req *http.Request, name string) string {
	return strings.Trim(req.FormValue(name), " ")
}

func (h *gameHandler) extractStringsParameter(req *http.Request, name string) []string {
	strValue := h.extractStringParameter(req, name)
	return strings.Split(strValue, ",")
}

func (h *gameHandler) extractBoolParameter(req *http.Request, name string) bool {
	strValue := h.extractStringParameter(req, name)
	strValue = strings.ToLower(strValue)
	return strValue == "true" || strValue == "1"
}

func (h *gameHandler) extractIntParameter(req *http.Request, name string) int {
	strValue := h.extractStringParameter(req, name)
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return 0
}

func (h *gameHandler) extractGameIdFromPath(req *http.Request) (model.GameId, error) {
	params := httprouter.ParamsFromContext(req.Context())
	strValue := params.ByName("game_id")
	if value, err := strconv.ParseInt(strValue, 10, 64); err == nil {
		return model.GameId(value), nil
	}
	return 0, model.ErrInvalidGameId
}

// //////////////////////////////////////////////////
// create

func (h *gameHandler) handleCreateGame(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var settings model.GameSettings
	var game *model.Game
	var err error

	//
	// decode request
	//

	settings = model.GameSettings{
		Seed:       time.Now().UnixMilli(),
		NbQuestion: h.extractIntParameter(req, "nb_question"),
		NbAnswer:   h.extractIntParameter(req, "nb_answer"),
		NbPlayer:   h.extractIntParameter(req, "nb_player"),
		Sources: util.Filter(
			util.Convert(
				h.extractStringsParameter(req, "sources"),
				model.ToSource,
			),
			func(s model.Source) bool { return s != "" },
		),
	}
	// CLEAN
	if len(settings.Sources) == 0 {
		h.logger.Warn("[api] missing sources >>> FALLBACK to legacy")
		settings.Sources = append(settings.Sources, model.Source_Legacy)
	}

	err = settings.Validate()
	if err != nil {
		goto encode_error
	}
	h.logger.Info(fmt.Sprintf("[api] create game with %d question(s), %d answer(s), %d player(s) and %d sources: %#v", settings.NbQuestion, settings.NbAnswer, settings.NbPlayer, len(settings.Sources), settings.Sources))

	//
	// execute
	//

	game, err = h.service.CreateGame(ctx, settings)
	if err != nil {
		goto encode_error
	}
	if game == nil {
		err = model.ErrGameNotFound
		goto encode_error
	}

	//
	// encode response
	//

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	err = json.NewEncoder(resp).Encode(h.toJsonGame(game))
	if err != nil {
		goto encode_error
	}
	return

encode_error:
	// TODO status code
	h.encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// retrieve

func (h *gameHandler) handleRetrieveGame(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var gameId model.GameId
	var game *model.Game
	var err error

	switch {
	default:

		//
		// decode request
		//

		gameId, err = h.extractGameIdFromPath(req)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] retrieve game %d", gameId))

		//
		// execute
		//

		game, err = h.service.RetrieveGame(ctx, gameId)
		if err != nil {
			break
		}
		if game == nil {
			err = model.ErrGameNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(h.toJsonGame(game))
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	h.encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// update

func (h *gameHandler) handleUpdateGame(resp http.ResponseWriter, req *http.Request) {

	// ctx := req.Context()

	var gameId model.GameId
	var game *model.Game
	var err error

	switch {
	default:

		//
		// decode request
		//

		gameId, err = h.extractGameIdFromPath(req)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] update game %d", gameId))

		//
		// execute
		//

		// TODO
		err = model.ErrNotImplemented
		if err != nil {
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(h.toJsonGame(game))
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	h.encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// delete

func (h *gameHandler) handleDeleteGame(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var gameId model.GameId
	var err error

	switch {
	default:

		//
		// decode request
		//

		gameId, err = h.extractGameIdFromPath(req)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] delete game %d", gameId))

		//
		// execute
		//

		err = h.service.DeleteGame(ctx, gameId)
		if err != nil {
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(h.toJsonSuccess())
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	h.encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// encode

func (h *gameHandler) toJsonSuccess() *JsonGameResponse {
	return &JsonGameResponse{
		Success: true,
	}
}

func (h *gameHandler) toJsonGame(game *model.Game) *JsonGameResponse {
	return &JsonGameResponse{
		Success: true,
		Game: &JsonGame{
			Id:        int64(game.Id),
			Settings:  h.toJsonGameSettings(game.Settings),
			Players:   util.Convert(game.Players, h.toJsonPlayer),
			Questions: util.Convert(game.Questions, h.toJsonQuestion),
		},
	}
}

func (h *gameHandler) toJsonGameSettings(settings *model.GameSettings) *JsonGameSettings {
	return &JsonGameSettings{
		Seed:       settings.Seed,
		NbQuestion: settings.NbQuestion,
		NbAnswer:   settings.NbAnswer,
		NbPlayer:   settings.NbPlayer,
		Sources:    util.Convert(settings.Sources, model.Source.String),
	}
}

func (h *gameHandler) toJsonPlayer(player *model.GamePlayer) *JsonPlayer {
	return &JsonPlayer{
		Id:     int64(player.Id),
		Name:   player.Name,
		Active: player.Active,
		Score:  player.Score,
	}
}

func (h *gameHandler) toJsonQuestion(question *model.GameQuestion) *JsonQuestion {
	return &JsonQuestion{
		Id:      int64(question.Id),
		Theme:   h.toJsonTheme(question.Theme),
		Music:   h.toJsonMusic(question.Music),
		Answers: util.Convert(question.Answers, h.toJsonAnswer),
	}
}

func (h *gameHandler) toJsonTheme(theme *model.GameTheme) JsonTheme {
	return JsonTheme{
		Id:    theme.Id,
		Title: theme.Title,
	}
}

func (h *gameHandler) toJsonMusic(music *model.Music) JsonMusic {
	return JsonMusic{
		Id:       int64(music.Id),
		DeezerId: int64(music.DeezerMusicId),
		Name:     music.Name,
		Mp3Url:   music.Mp3Url,
		Artist:   h.toJsonArtist(music.Artist),
		Album:    h.toJsonAlbum(music.Album),
		Genre:    h.toJsonGenre(music.Genre),
	}
}

func (h *gameHandler) toJsonArtist(artist *model.MusicArtist) *JsonArtist {
	if artist == nil {
		return nil
	}
	return &JsonArtist{
		Id:       int64(artist.Id),
		DeezerId: int64(artist.DeezerArtistId),
		Name:     artist.Name,
		ImgUrl:   artist.ImgUrl,
	}
}

func (h *gameHandler) toJsonAlbum(album *model.MusicAlbum) *JsonAlbum {
	if album == nil {
		return nil
	}
	return &JsonAlbum{
		Id:       int64(album.Id),
		DeezerId: int64(album.DeezerAlbumId),
		Name:     album.Name,
		ImgUrl:   album.ImgUrl,
	}
}

func (h *gameHandler) toJsonGenre(genre *model.MusicGenre) *JsonGenre {
	if genre == nil {
		return nil
	}
	return &JsonGenre{
		Id:       int64(genre.Id),
		DeezerId: int64(genre.DeezerGenreId),
		Name:     genre.Name,
		ImgUrl:   genre.ImgUrl,
	}
}

func (h *gameHandler) toJsonAnswer(answer *model.GameAnswer) JsonAnswer {
	return JsonAnswer{
		Id:      int64(answer.Id),
		Text:    answer.Text,
		Hint:    answer.Hint,
		Correct: answer.Correct,
	}
}

type JsonGameResponse struct {
	Success bool       `json:"success,omitempty"`
	Error   *JsonError `json:"error,omitempty"`
	Game    *JsonGame  `json:"game,omitempty"`
}

type JsonGame struct {
	Id        int64             `json:"id,omitempty"`
	Settings  *JsonGameSettings `json:"settings,omitempty"`
	Players   []*JsonPlayer     `json:"players,omitempty"`
	Questions []*JsonQuestion   `json:"questions,omitempty"`
}

type JsonGameSettings struct {
	Seed       int64    `json:"seed,omitempty"`
	NbQuestion int      `json:"nbQuestion,omitempty"`
	NbAnswer   int      `json:"nbAnswer,omitempty"`
	NbPlayer   int      `json:"nbPlayer,omitempty"`
	Sources    []string `json:"sources,omitempty"`
}

type JsonPlayer struct {
	Id     int64  `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Active bool   `json:"active,omitempty"`
	Score  int    `json:"score,omitempty"`
}

type JsonQuestion struct {
	Id      int64        `json:"id"`
	Theme   JsonTheme    `json:"theme"`
	Music   JsonMusic    `json:"music"`
	Answers []JsonAnswer `json:"answers,omitempty"`
}

type JsonTheme struct {
	Id    int64  `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type JsonMusic struct {
	Id       int64       `json:"id,omitempty"`
	DeezerId int64       `json:"deezerId,omitempty"`
	Name     string      `json:"name,omitempty"`
	Mp3Url   string      `json:"mp3Url,omitempty"`
	Artist   *JsonArtist `json:"artist,omitempty"`
	Album    *JsonAlbum  `json:"album,omitempty"`
	Genre    *JsonGenre  `json:"genre,omitempty"`
}

type JsonArtist struct {
	Id       int64  `json:"id,omitempty"`
	DeezerId int64  `json:"deezerId,omitempty"`
	Name     string `json:"name,omitempty"`
	ImgUrl   string `json:"imgUrl,omitempty"`
}

type JsonAlbum struct {
	Id       int64  `json:"id,omitempty"`
	DeezerId int64  `json:"deezerId,omitempty"`
	Name     string `json:"name,omitempty"`
	ImgUrl   string `json:"imgUrl,omitempty"`
}

type JsonGenre struct {
	Id       int64  `json:"id,omitempty"`
	DeezerId int64  `json:"deezerId,omitempty"`
	Name     string `json:"name,omitempty"`
	ImgUrl   string `json:"imgUrl,omitempty"`
}

type JsonAnswer struct {
	Id      int64  `json:"id"`
	Text    string `json:"text"`
	Hint    string `json:"hint,omitempty"`
	Correct bool   `json:"correct,omitempty"`
}

// //////////////////////////////////////////////////
// encode error

func (h *gameHandler) encodeError(resp http.ResponseWriter, statusCode int, message string) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)

	// try to encode error >>> no need to check error at encoding
	json.NewEncoder(resp).Encode(h.toJsonError(statusCode, message))
}

func (h *gameHandler) toJsonError(code int, message string) *JsonGameResponse {
	return &JsonGameResponse{
		Success: false,
		Error: &JsonError{
			Code:    code,
			Message: message,
		},
	}
}

type JsonError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
