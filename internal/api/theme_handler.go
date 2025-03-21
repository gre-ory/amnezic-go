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

func NewThemehandler(logger *zap.Logger, themeService service.ThemeService, musicService service.MusicService, sessionService service.SessionService) Handler {
	return &themeHandler{
		logger:         logger,
		themeService:   themeService,
		musicService:   musicService,
		sessionService: sessionService,
	}
}

type themeHandler struct {
	logger         *zap.Logger
	themeService   service.ThemeService
	musicService   service.MusicService
	sessionService service.SessionService
}

// //////////////////////////////////////////////////
// register

func (h *themeHandler) RegisterRoutes(router *httprouter.Router) {

	router.HandlerFunc(http.MethodGet, "/api/theme", h.handleListTheme)
	router.HandlerFunc(http.MethodGet, "/api/theme/:theme_id", h.handleRetrieveTheme)

	withThemePermission := WithPermission(h.logger, h.sessionService, model.Permission_Theme)

	router.HandlerFunc(http.MethodPut, "/api/theme/new", withThemePermission(h.handleCreateTheme))
	router.HandlerFunc(http.MethodPost, "/api/theme/:theme_id", withThemePermission(h.handleUpdateTheme))
	router.HandlerFunc(http.MethodDelete, "/api/theme/:theme_id", withThemePermission(h.handleDeleteTheme))
	router.HandlerFunc(http.MethodPut, "/api/theme-question/:theme_id/new", withThemePermission(h.handleAddQuestion))
	router.HandlerFunc(http.MethodPost, "/api/theme-question/:theme_id/:question_id", withThemePermission(h.handleUpdateQuestion))
	router.HandlerFunc(http.MethodDelete, "/api/theme-question/:theme_id/:question_id", withThemePermission(h.handleRemoveQuestion))
}

// //////////////////////////////////////////////////
// list

func (h *themeHandler) handleListTheme(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var themes []*model.ThemeInfo
	var err error

	switch {
	default:

		h.logger.Info("[api] list themes")

		//
		// execute
		//

		themes, err = h.themeService.ListThemes(ctx)
		if err != nil {
			break
		}

		//
		// encode response
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonThemesResponse(themes))
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

		theme, err = extractThemeFromBody(req, h.logger)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] create theme: %#v", theme))

		//
		// execute
		//

		theme, err = h.themeService.CreateTheme(ctx, theme)
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

		theme, err = h.themeService.RetrieveTheme(ctx, themeId)
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
		theme, err = extractThemeFromBody(req, h.logger)
		if err != nil {
			break
		}
		if themeId != theme.Id {
			err = model.ErrInvalidThemeId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] update theme %d: %#v", themeId, theme))

		//
		// execute
		//

		theme, err = h.themeService.UpdateTheme(ctx, theme)
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

		err = h.themeService.DeleteTheme(ctx, themeId)
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
// add question

func (h *themeHandler) handleAddQuestion(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var themeId model.ThemeId
	var deezerId model.DeezerMusicId
	var music *model.Music
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

		deezerId = model.DeezerMusicId(toInt64(extractParameter(req, "deezer_id")))
		if deezerId == 0 {
			err = model.ErrInvalidDeezerId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] add question from deezer music %d to theme %d", deezerId, themeId))

		//
		// execute
		//

		music, err = h.musicService.AddDeezerMusic(ctx, deezerId)
		if err != nil {
			break
		}
		if music == nil {
			err = model.ErrMusicNotFound
			break
		}

		//
		// execute
		//

		question := music.ToThemeQuestion(themeId)
		theme, err = h.themeService.AddQuestion(ctx, question)
		if err != nil {
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
// update question

func (h *themeHandler) handleUpdateQuestion(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var themeId model.ThemeId
	var questionId model.ThemeQuestionId
	var question *model.ThemeQuestion
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
		questionId = model.ThemeQuestionId(toInt64(extractPathParameter(req, "question_id")))
		if questionId == 0 {
			err = model.ErrInvalidThemeQuestionId
			break
		}
		question, err = extractThemeQuestionFromBody(req, h.logger, themeId)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] update question %d from theme %d", questionId, themeId), zap.Object("question", question))

		//
		// execute
		//

		theme, err = h.themeService.UpdateQuestion(ctx, question)
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
// remove question

func (h *themeHandler) handleRemoveQuestion(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var themeId model.ThemeId
	var questionId model.ThemeQuestionId
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
		questionId = model.ThemeQuestionId(toInt64(extractPathParameter(req, "question_id")))
		if questionId == 0 {
			err = model.ErrInvalidThemeQuestionId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] remove question %d from theme %d", questionId, themeId))

		//
		// execute
		//

		theme, err = h.themeService.RemoveQuestion(ctx, themeId, questionId)
		if err != nil {
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
// decode

func extractThemeFromBody(req *http.Request, logger *zap.Logger) (*model.Theme, error) {
	var jsonBody JsonThemeBody
	jsonErr := json.NewDecoder(req.Body).Decode(&jsonBody)
	switch {
	case jsonErr == io.EOF:
		logger.Info("failed to decode theme body: EOF")
		return nil, model.ErrInvalidBody
	case jsonErr != nil:
		logger.Info("failed to decode theme body", zap.Error(jsonErr))
		return nil, model.ErrInvalidBody
	}

	return toTheme(jsonBody.Theme), nil
}

func extractThemeQuestionFromBody(req *http.Request, logger *zap.Logger, themeId model.ThemeId) (*model.ThemeQuestion, error) {
	var jsonBody JsonThemeQuestionBody
	jsonErr := json.NewDecoder(req.Body).Decode(&jsonBody)
	switch {
	case jsonErr == io.EOF:
		logger.Info("failed to decode theme body: EOF")
		return nil, model.ErrInvalidBody
	case jsonErr != nil:
		logger.Info("failed to decode theme question body", zap.Error(jsonErr))
		return nil, model.ErrInvalidBody
	}

	return toThemeQuestion(themeId, jsonBody.Question), nil
}

func toTheme(jsonTheme *JsonTheme) *model.Theme {
	theme := &model.Theme{
		Id:     model.ThemeId(jsonTheme.Id),
		Title:  jsonTheme.Title,
		ImgUrl: jsonTheme.ImgUrl,
		Labels: jsonTheme.Labels,
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

type JsonThemeQuestionBody struct {
	Question *JsonThemeQuestion `json:"question,omitempty"`
}

// //////////////////////////////////////////////////
// encode

func toJsonThemesResponse(themes []*model.ThemeInfo) *JsonThemesResponse {
	return &JsonThemesResponse{
		Success: true,
		Themes:  util.Convert(themes, toJsonThemeInfo),
	}
}

func toJsonThemeInfo(theme *model.ThemeInfo) *JsonThemeInfo {
	return &JsonThemeInfo{
		Id:         int64(theme.Id),
		Title:      theme.Title,
		ImgUrl:     theme.ImgUrl,
		Labels:     theme.Labels,
		NbQuestion: theme.NbQuestion,
	}
}

func toJsonThemeResponse(theme *model.Theme) *JsonThemeResponse {
	return &JsonThemeResponse{
		Success: true,
		Theme:   toJsonTheme(theme),
	}
}

func toJsonTheme(theme *model.Theme) *JsonTheme {
	jsonTheme := &JsonTheme{
		Id:     int64(theme.Id),
		Title:  theme.Title,
		ImgUrl: theme.ImgUrl,
		Labels: theme.Labels,
	}

	if theme.Questions != nil {
		jsonTheme.Questions = util.Convert(theme.Questions, toJsonThemeQuestion)
	}

	return jsonTheme
}

func toJsonThemeQuestion(question *model.ThemeQuestion) *JsonThemeQuestion {
	jsonQuestion := &JsonThemeQuestion{
		Id:   int64(question.Id),
		Text: question.Text,
		Hint: question.Hint,
	}

	if question.Music != nil {
		jsonQuestion.Music = toJsonMusic(question.Music)
	}

	if question.Theme != nil {
		jsonQuestion.Theme = toJsonTheme(question.Theme)
	}

	return jsonQuestion
}

type JsonThemesResponse struct {
	Success bool             `json:"success,omitempty"`
	Themes  []*JsonThemeInfo `json:"themes"`
}

type JsonThemeInfo struct {
	Id         int64             `json:"id,omitempty"`
	Title      string            `json:"title,omitempty"`
	ImgUrl     string            `json:"imgUrl,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
	NbQuestion int               `json:"nbQuestion,omitempty"`
}

type JsonThemeResponse struct {
	Success bool       `json:"success,omitempty"`
	Theme   *JsonTheme `json:"theme,omitempty"`
}

type JsonTheme struct {
	Id        int64                `json:"id,omitempty"`
	Title     string               `json:"title,omitempty"`
	ImgUrl    string               `json:"imgUrl,omitempty"`
	Labels    map[string]string    `json:"labels,omitempty"`
	Questions []*JsonThemeQuestion `json:"questions,omitempty"`
}

type JsonThemeQuestion struct {
	Id    int64      `json:"id,omitempty"`
	Text  string     `json:"text,omitempty"`
	Hint  string     `json:"hint,omitempty"`
	Theme *JsonTheme `json:"theme,omitempty"`
	Music *JsonMusic `json:"music,omitempty"`
}
