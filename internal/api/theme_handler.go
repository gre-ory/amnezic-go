package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/service"
	"github.com/gre-ory/amnezic-go/internal/util"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// theme handler

func NewThemehandler(logger *zap.Logger, service service.ThemeService) Handler {
	return &themeHandler{
		logger:  logger,
		service: service,
	}
}

type themeHandler struct {
	logger  *zap.Logger
	service service.ThemeService
}

// //////////////////////////////////////////////////
// register

func (h *themeHandler) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPut, "/api/theme/new", h.handleCreateTheme)
	router.HandlerFunc(http.MethodGet, "/api/theme/:theme_id", h.handleRetrieveTheme)
	router.HandlerFunc(http.MethodPost, "/api/theme/:theme_id", h.handleUpdateTheme)
	router.HandlerFunc(http.MethodDelete, "/api/theme/:theme_id", h.handleDeleteTheme)
}

// //////////////////////////////////////////////////
// create

func (h *themeHandler) handleCreateTheme(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var theme *model.Theme
	var err error

	switch {
	default:

		//
		// decode request
		//

		theme, err = extractThemeFromBody(req)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] create theme: %#v", theme))

		//
		// execute
		//

		theme, err = h.service.CreateTheme(ctx, theme)
		if err != nil {
			break
		}
		if theme == nil {
			err = model.ErrThemeNotFound
			break
		}

		//
		// encode response
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonThemeResponse(theme))
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

func (h *themeHandler) handleRetrieveTheme(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var themeId model.ThemeId
	var theme *model.Theme
	var err error

	switch {
	default:

		//
		// decode request
		//

		themeId = model.ThemeId(toInt64(extractPathParameter(req, "theme_id")))
		if themeId == 0 {
			err = model.ErrInvalidThemeId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] retrieve theme %d", themeId))

		//
		// execute
		//

		theme, err = h.service.RetrieveTheme(ctx, themeId)
		if err != nil {
			break
		}
		if theme == nil {
			err = model.ErrThemeNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonThemeResponse(theme))
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

func (h *themeHandler) handleUpdateTheme(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var themeId model.ThemeId
	var theme *model.Theme
	var err error

	switch {
	default:

		//
		// decode request
		//

		themeId = model.ThemeId(toInt64(extractPathParameter(req, "theme_id")))
		if themeId == 0 {
			err = model.ErrInvalidThemeId
			break
		}
		theme, err = extractThemeFromBody(req)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] update theme %d: %#v", themeId, theme))

		//
		// execute
		//

		theme, err = h.service.UpdateTheme(ctx, theme)
		if err != nil {
			break
		}
		if theme == nil {
			err = model.ErrThemeNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonThemeResponse(theme))
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

func (h *themeHandler) handleDeleteTheme(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var themeId model.ThemeId
	var err error

	switch {
	default:

		//
		// decode request
		//

		themeId = model.ThemeId(toInt64(extractPathParameter(req, "theme_id")))
		if themeId == 0 {
			err = model.ErrInvalidThemeId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] delete theme %d", themeId))

		//
		// execute
		//

		err = h.service.DeleteTheme(ctx, themeId)
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
// decode

func extractThemeFromBody(req *http.Request) (*model.Theme, error) {
	var jsonBody JsonThemeBody
	jsonErr := json.NewDecoder(req.Body).Decode(&jsonBody)
	switch {
	case jsonErr == io.EOF:
		return nil, model.ErrInvalidBody
	case jsonErr != nil:
		return nil, model.ErrInvalidBody
	}

	return toTheme(jsonBody.Theme), nil
}

func toTheme(jsonTheme *JsonTheme) *model.Theme {
	theme := &model.Theme{
		Id:     model.ThemeId(jsonTheme.Id),
		Title:  jsonTheme.Title,
		ImgUrl: jsonTheme.ImgUrl,
	}

	theme.Questions = util.Convert(jsonTheme.Questions, func(jsonQuestion *JsonThemeQuestion) *model.ThemeQuestion {
		return toThemeQuestion(theme.Id, jsonQuestion)
	})

	return theme
}

func toThemeQuestion(themeId model.ThemeId, jsonQuestion *JsonThemeQuestion) *model.ThemeQuestion {
	question := &model.ThemeQuestion{
		Id:      model.ThemeQuestionId(jsonQuestion.Id),
		ThemeId: themeId,
		Text:    jsonQuestion.Text,
		Hint:    jsonQuestion.Hint,
	}

	if jsonQuestion.Music != nil {
		question.MusicId = model.MusicId(jsonQuestion.Music.Id)
	}

	return question
}

type JsonThemeBody struct {
	Theme *JsonTheme `json:"theme,omitempty"`
}

// //////////////////////////////////////////////////
// encode

func toJsonThemeResponse(theme *model.Theme) *JsonThemeResponse {
	return &JsonThemeResponse{
		Success: true,
		Theme:   toJsonTheme(theme),
	}
}

func toJsonTheme(theme *model.Theme) *JsonTheme {
	return &JsonTheme{
		Id:        int64(theme.Id),
		Title:     theme.Title,
		ImgUrl:    theme.ImgUrl,
		Questions: util.Convert(theme.Questions, toJsonThemeQuestion),
	}
}

func toJsonThemeQuestion(question *model.ThemeQuestion) *JsonThemeQuestion {
	jsonQuestion := &JsonThemeQuestion{
		Id:   int64(question.Id),
		Text: question.Text,
		Hint: question.Hint,
	}

	if question.Music != nil {
		jsonQuestion.Music = toJsonMusic(question.Music)
	} else {
		jsonQuestion.Music = &JsonMusic{
			Id: int64(question.MusicId),
		}
	}

	return jsonQuestion
}

type JsonThemeResponse struct {
	Success bool       `json:"success,omitempty"`
	Theme   *JsonTheme `json:"theme,omitempty"`
}

type JsonTheme struct {
	Id        int64                `json:"id,omitempty"`
	Title     string               `json:"title,omitempty"`
	ImgUrl    string               `json:"imgUrl,omitempty"`
	Questions []*JsonThemeQuestion `json:"questions,omitempty"`
}

type JsonThemeQuestion struct {
	Id    int64      `json:"id,omitempty"`
	Text  string     `json:"text,omitempty"`
	Hint  string     `json:"hint,omitempty"`
	Music *JsonMusic `json:"music,omitempty"`
}
