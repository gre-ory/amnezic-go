package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	router.HandlerFunc(http.MethodPut, "/game/new", h.handleCreateGame)
	router.HandlerFunc(http.MethodGet, "/game/:game_id", h.handleRetrieveGame)
	router.HandlerFunc(http.MethodPost, "/game/:game_id", h.handleUpdateGame)
	router.HandlerFunc(http.MethodDelete, "/game/:game_id", h.handleDeleteGame)
}

// //////////////////////////////////////////////////
// decode

func (h *gameHandler) extractIntParameter(req *http.Request, name string) int {
	strValue := req.FormValue(name)
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
		UseLegacy:  true,
		NbQuestion: h.extractIntParameter(req, "nb_question"),
		NbAnswer:   h.extractIntParameter(req, "nb_answer"),
		NbPlayer:   h.extractIntParameter(req, "nb_player"),
	}
	err = settings.Validate()
	if err != nil {
		goto encode_error
	}
	h.logger.Info(fmt.Sprintf("[api] create game with %d question(s), %d answer(s) and %d player(s)", settings.NbQuestion, settings.NbAnswer, settings.NbPlayer))

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

	//
	// decode request
	//

	gameId, err = h.extractGameIdFromPath(req)
	if err != nil {
		goto encode_error
	}
	h.logger.Info(fmt.Sprintf("[api] retrieve game %d", gameId))

	//
	// execute
	//

	game, err = h.service.RetrieveGame(ctx, gameId)
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
// update

func (h *gameHandler) handleUpdateGame(resp http.ResponseWriter, req *http.Request) {

	// ctx := req.Context()

	var gameId model.GameId
	var game *model.Game
	var err error

	//
	// decode request
	//

	gameId, err = h.extractGameIdFromPath(req)
	if err != nil {
		goto encode_error
	}
	h.logger.Info(fmt.Sprintf("[api] update game %d", gameId))

	//
	// execute
	//

	// TODO
	err = model.ErrNotImplemented
	if err != nil {
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
// delete

func (h *gameHandler) handleDeleteGame(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var gameId model.GameId
	var err error

	//
	// decode request
	//

	gameId, err = h.extractGameIdFromPath(req)
	if err != nil {
		goto encode_error
	}
	h.logger.Info(fmt.Sprintf("[api] delete game %d", gameId))

	//
	// execute
	//

	err = h.service.DeleteGame(ctx, gameId)
	if err != nil {
		goto encode_error
	}

	//
	// encode response
	//

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	err = json.NewEncoder(resp).Encode(h.toJsonSuccess())
	if err != nil {
		goto encode_error
	}
	return

encode_error:
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
			Players:   util.Convert(game.Players, h.toJsonPlayer),
			Questions: util.Convert(game.Questions, h.toJsonQuestion),
		},
	}
}

func (h *gameHandler) toJsonPlayer(player *model.Player) *JsonPlayer {
	return &JsonPlayer{
		Id:     int64(player.Id),
		Name:   player.Name,
		Active: player.Active,
		Score:  player.Score,
	}
}

func (h *gameHandler) toJsonQuestion(question *model.Question) *JsonQuestion {
	return &JsonQuestion{
		Id:      int64(question.Id),
		Theme:   h.toJsonTheme(question.Theme),
		Music:   h.toJsonMusic(question.Music),
		Answers: util.Convert(question.Answers, h.toJsonAnswer),
	}
}

func (h *gameHandler) toJsonTheme(theme model.Theme) JsonTheme {
	return JsonTheme{
		Id:    theme.Id,
		Title: theme.Title,
	}
}

func (h *gameHandler) toJsonMusic(music model.Music) JsonMusic {
	return JsonMusic{
		Id:     music.Id,
		Name:   music.Name,
		Mp3Url: music.Mp3Url,
		Artist: h.toJsonArtist(music.Artist),
		Album:  h.toJsonAlbum(music.Album),
		Genre:  h.toJsonGenre(music.Genre),
	}
}

func (h *gameHandler) toJsonArtist(artist *model.Artist) *JsonArtist {
	if artist == nil {
		return nil
	}
	return &JsonArtist{
		Id:     artist.Id,
		Name:   artist.Name,
		ImgUrl: artist.ImgUrl,
	}
}

func (h *gameHandler) toJsonAlbum(album *model.Album) *JsonAlbum {
	if album == nil {
		return nil
	}
	return &JsonAlbum{
		Id:     album.Id,
		Name:   album.Name,
		ImgUrl: album.ImgUrl,
	}
}

func (h *gameHandler) toJsonGenre(genre *model.Genre) *JsonGenre {
	if genre == nil {
		return nil
	}
	return &JsonGenre{
		Id:     genre.Id,
		Name:   genre.Name,
		ImgUrl: genre.ImgUrl,
	}
}

func (h *gameHandler) toJsonAnswer(answer *model.Answer) JsonAnswer {
	return JsonAnswer{
		Id:    int64(answer.Id),
		Text:  answer.Text,
		Hint:  answer.Hint,
		Valid: answer.Correct,
	}
}

type JsonGameResponse struct {
	Success bool       `json:"success,omitempty"`
	Error   *JsonError `json:"error,omitempty"`
	Game    *JsonGame  `json:"game,omitempty"`
}

type JsonGame struct {
	Id        int64           `json:"id,omitempty"`
	Players   []*JsonPlayer   `json:"players,omitempty"`
	Questions []*JsonQuestion `json:"questions,omitempty"`
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
	Id     int64       `json:"id,omitempty"`
	Name   string      `json:"name,omitempty"`
	Mp3Url string      `json:"mp3Url,omitempty"`
	Artist *JsonArtist `json:"artist,omitempty"`
	Album  *JsonAlbum  `json:"album,omitempty"`
	Genre  *JsonGenre  `json:"genre,omitempty"`
}

type JsonArtist struct {
	Id     int64  `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	ImgUrl string `json:"imgUrl,omitempty"`
}

type JsonAlbum struct {
	Id     int64  `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	ImgUrl string `json:"imgUrl,omitempty"`
}

type JsonGenre struct {
	Id     int64  `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	ImgUrl string `json:"imgUrl,omitempty"`
}

type JsonAnswer struct {
	Id    int64  `json:"id"`
	Text  string `json:"text"`
	Hint  string `json:"hint,omitempty"`
	Valid bool   `json:"valid,omitempty"`
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
