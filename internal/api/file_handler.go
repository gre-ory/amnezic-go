package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/service"
	"github.com/gre-ory/amnezic-go/internal/util"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// file handler

func NewFilehandler(logger *zap.Logger, musicFilter *model.FileFilter, imageFilter *model.FileFilter, fileService service.FileService, sessionService service.SessionService) Handler {
	return &fileHandler{
		logger:         logger,
		musicFilter:    musicFilter,
		imageFilter:    imageFilter,
		fileService:    fileService,
		sessionService: sessionService,
	}
}

type fileHandler struct {
	logger         *zap.Logger
	musicFilter    *model.FileFilter
	imageFilter    *model.FileFilter
	fileService    service.FileService
	sessionService service.SessionService
}

// //////////////////////////////////////////////////
// register

func (h *fileHandler) RegisterRoutes(router *httprouter.Router) {

	router.ServeFiles("/static/music/*filepath", NewFilteredDirectory(h.musicFilter))
	router.ServeFiles("/static/image/*filepath", NewFilteredDirectory(h.imageFilter))

	withSessionPermission := WithPermission(h.logger, h.sessionService, model.Permission_File)

	router.HandlerFunc(http.MethodGet, "/file/music", withSessionPermission(h.handleListMusic))
	router.HandlerFunc(http.MethodGet, "/file/image", withSessionPermission(h.handleListImage))
}

// //////////////////////////////////////////////////
// serve

type FilteredDirectory struct {
	Directory http.FileSystem
	Filter    *model.FileFilter
}

func NewFilteredDirectory(filter *model.FileFilter) *FilteredDirectory {
	return &FilteredDirectory{
		Directory: http.Dir(filter.Directory),
		Filter:    filter,
	}
}

func (fs *FilteredDirectory) Open(name string) (http.File, error) {
	if !fs.Filter.MatchExtension(name) {
		return nil, os.ErrNotExist
	}
	return fs.Directory.Open(name)
}

// //////////////////////////////////////////////////
// list

func (h *fileHandler) handleListMusic(resp http.ResponseWriter, req *http.Request) {
	h.handleList(resp, req, h.musicFilter)
}

func (h *fileHandler) handleListImage(resp http.ResponseWriter, req *http.Request) {
	h.handleList(resp, req, h.imageFilter)
}

func (h *fileHandler) handleList(resp http.ResponseWriter, req *http.Request, filter *model.FileFilter) {
	defer onPanic(resp)()

	ctx := req.Context()

	var urls []model.Url
	var err error

	switch {
	default:

		h.logger.Info("[api] list files")

		//
		// execute
		//

		urls, err = h.fileService.List(ctx, filter)
		if err != nil {
			break
		}

		//
		// encode response
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonUrlsResponse(urls))
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

func toJsonUrlsResponse(urls []model.Url) *JsonUrlsResponse {
	return &JsonUrlsResponse{
		Success: true,
		Urls:    util.Convert(urls, util.ToStr[model.Url]),
	}
}

type JsonUrlsResponse struct {
	Success bool     `json:"success,omitempty"`
	Urls    []string `json:"urls"`
}
