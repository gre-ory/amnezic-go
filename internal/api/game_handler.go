package api

import (
	"encoding/json"
	"fmt"
	"net/http"
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
// create

func (h *gameHandler) handleCreateGame(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var settings model.GameSettings
	var game *model.Game
	var err error

	switch {
	default:

		//
		// decode request
		//

		settings = model.GameSettings{
			Seed:       time.Now().UnixMilli(),
			NbQuestion: toInt(extractParameter(req, "nb_question")),
			NbAnswer:   toInt(extractParameter(req, "nb_answer")),
			NbPlayer:   toInt(extractParameter(req, "nb_player")),
			Sources: util.Filter(
				util.Convert(
					toStrings(extractParameter(req, "sources")),
					model.ToSource,
				),
				func(s model.Source) bool { return s != "" },
			),
		}
		// CLEAN
		if len(settings.Sources) == 0 {
			h.logger.Info("[api] missing sources >>> FALLBACK to store")
			settings.Sources = append(settings.Sources, model.Source_Store)
		}

		err = settings.Validate()
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] create game with %d question(s), %d answer(s), %d player(s) and %d sources: %#v", settings.NbQuestion, settings.NbAnswer, settings.NbPlayer, len(settings.Sources), settings.Sources))

		//
		// execute
		//

		game, err = h.service.CreateGame(ctx, settings)
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
		err = json.NewEncoder(resp).Encode(toJsonGameResponse(game))
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	encodeError(resp, http.StatusBadRequest, err.Error())
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

		gameId = model.GameId(toInt64(extractPathParameter(req, "game_id")))
		if gameId == 0 {
			err = model.ErrInvalidGameId
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
		err = json.NewEncoder(resp).Encode(toJsonGameResponse(game))
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	encodeError(resp, http.StatusBadRequest, err.Error())
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

		gameId = model.GameId(toInt64(extractPathParameter(req, "game_id")))
		if gameId == 0 {
			err = model.ErrInvalidGameId
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
		err = json.NewEncoder(resp).Encode(toJsonGameResponse(game))
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	encodeError(resp, http.StatusBadRequest, err.Error())
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

		gameId = model.GameId(toInt64(extractPathParameter(req, "game_id")))
		if gameId == 0 {
			err = model.ErrInvalidGameId
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
		err = json.NewEncoder(resp).Encode(toJsonSuccess())
		if err != nil {
			break
		}
		return
	}

	//
	// encode error
	//

	// TODO status code
	encodeError(resp, http.StatusBadRequest, err.Error())
}

// //////////////////////////////////////////////////
// encode

func toJsonGameResponse(game *model.Game) *JsonGameResponse {
	return &JsonGameResponse{
		Success: true,
		Game:    toJsonGame(game),
	}
}

func toJsonGame(game *model.Game) *JsonGame {
	return &JsonGame{
		Id:        int64(game.Id),
		Settings:  toJsonGameSettings(game.Settings),
		Players:   util.Convert(game.Players, toJsonGamePlayer),
		Questions: util.Convert(game.Questions, toJsonGameQuestion),
	}
}

func toJsonGameSettings(settings *model.GameSettings) *JsonGameSettings {
	return &JsonGameSettings{
		Seed:       settings.Seed,
		NbQuestion: settings.NbQuestion,
		NbAnswer:   settings.NbAnswer,
		NbPlayer:   settings.NbPlayer,
		Sources:    util.Convert(settings.Sources, model.Source.String),
	}
}

func toJsonGamePlayer(player *model.GamePlayer) *JsonGamePlayer {
	return &JsonGamePlayer{
		Id:     int64(player.Id),
		Name:   player.Name,
		Active: player.Active,
		Score:  player.Score,
	}
}

func toJsonGameQuestion(question *model.GameQuestion) *JsonGameQuestion {
	return &JsonGameQuestion{
		Id:      int64(question.Id),
		Theme:   toJsonGameTheme(question.Theme),
		Music:   toJsonMusic(question.Music),
		Answers: util.Convert(question.Answers, toJsonGameAnswer),
	}
}

func toJsonGameTheme(theme *model.GameTheme) *JsonGameTheme {
	return &JsonGameTheme{
		Id:    theme.Id,
		Title: theme.Title,
	}
}

func toJsonGameAnswer(answer *model.GameAnswer) *JsonGameAnswer {
	return &JsonGameAnswer{
		Id:      int64(answer.Id),
		Text:    answer.Text,
		Hint:    answer.Hint,
		Correct: answer.Correct,
	}
}

type JsonGameResponse struct {
	Success bool      `json:"success,omitempty"`
	Game    *JsonGame `json:"game,omitempty"`
}

type JsonGame struct {
	Id        int64               `json:"id,omitempty"`
	Settings  *JsonGameSettings   `json:"settings,omitempty"`
	Players   []*JsonGamePlayer   `json:"players,omitempty"`
	Questions []*JsonGameQuestion `json:"questions,omitempty"`
}

type JsonGameSettings struct {
	Seed       int64    `json:"seed,omitempty"`
	NbQuestion int      `json:"nbQuestion,omitempty"`
	NbAnswer   int      `json:"nbAnswer,omitempty"`
	NbPlayer   int      `json:"nbPlayer,omitempty"`
	Sources    []string `json:"sources,omitempty"`
}

type JsonGamePlayer struct {
	Id     int64  `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Active bool   `json:"active,omitempty"`
	Score  int    `json:"score,omitempty"`
}

type JsonGameQuestion struct {
	Id      int64             `json:"id"`
	Theme   *JsonGameTheme    `json:"theme"`
	Music   *JsonMusic        `json:"music"`
	Answers []*JsonGameAnswer `json:"answers,omitempty"`
}

type JsonGameTheme struct {
	Id    int64  `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type JsonGameAnswer struct {
	Id      int64  `json:"id"`
	Text    string `json:"text"`
	Hint    string `json:"hint,omitempty"`
	Correct bool   `json:"correct,omitempty"`
}
