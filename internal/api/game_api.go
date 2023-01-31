package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/service"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// game api

func NewGameHandler(service service.GameService) http.Handler {
	return &gameHandler{
		service: service,
	}
}

type gameHandler struct {
	service service.GameService
}

var (
	GamePathRegex   = regexp.MustCompile(`/game/`)
	GameIdPathRegex = regexp.MustCompile(`/game/(?P<id>\d+)`)
)

func (h *gameHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	ctx := context.Background()
	logger := zap.L()

	var game *model.Game
	var err error

	switch req.Method {
	case "GET":
		if values, ok := match(req, GameIdPathRegex); ok {
			game, err = h.service.RetrieveGame(ctx, logger, values["id"])
		} else if _, ok := match(req, GamePathRegex); ok {
			game, err = h.service.CreateGame(ctx, logger)
		} else {
			err = model.ErrNotImplemented
		}
	case "POST":
		err = model.ErrNotImplemented
	case "PUT":
		err = model.ErrNotImplemented
	case "DELETE":
		err = model.ErrNotImplemented
	default:
		err = model.ErrNotImplemented
	}

	if game != nil {
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonGame(game))
	}

	if err != nil {
		encodeError(resp, http.StatusBadRequest, err.Error())
	}
}

func match(req *http.Request, regex *regexp.Regexp) (map[string]int64, bool) {
	match := regex.FindStringSubmatch(req.URL.Path)
	if len(match) == 0 {
		return nil, false
	}
	result := make(map[string]int64, len(match))
	for i, name := range regex.SubexpNames() {
		fmt.Printf("regex #%d - %s=%s \n", i, name, match[i])
		if i != 0 && name != "" {
			value, _ := strconv.ParseInt(match[i], 10, 64)
			result[name] = value
		}
	}
	fmt.Printf("result=%v \n", result)
	return result, true
}

func encodeError(resp http.ResponseWriter, statusCode int, message string) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)

	// try to encode error >>> no need to check error at encoding
	json.NewEncoder(resp).Encode(toJsonError(statusCode, message))
}

func toJsonGame(game *model.Game) *JsonGame {
	return &JsonGame{
		Id:      game.Id,
		Players: adaptList(game.Players, toJsonPlayer),
	}
}

func adaptList[T any, J any](items []T, adapt func(item T) J) []J {
	result := make([]J, 0, len(items))
	for _, item := range items {
		result = append(result, adapt(item))
	}
	return result
}

func toJsonPlayer(player *model.Player) *JsonPlayer {
	return &JsonPlayer{
		Id:     player.Id,
		Name:   player.Name,
		Active: player.Active,
		Score:  player.Score,
	}
}

func toJsonError(code int, message string) *JsonErrorResponse {
	return &JsonErrorResponse{
		Error: JsonError{
			Code:    code,
			Message: message,
		},
	}
}

type JsonGame struct {
	Id      int64         `json:"id"`
	Players []*JsonPlayer `json:"players"`
}

type JsonPlayer struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
	Score  int    `json:"score"`
}

type JsonErrorResponse struct {
	Error JsonError `json:"error"`
}

type JsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
