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
// album handler

func NewAlbumhandler(logger *zap.Logger, service service.AlbumService, sessionService service.SessionService) Handler {
	return &albumHandler{
		logger:         logger,
		service:        service,
		sessionService: sessionService,
	}
}

type albumHandler struct {
	logger         *zap.Logger
	service        service.AlbumService
	sessionService service.SessionService
}

// //////////////////////////////////////////////////
// register

func (h *albumHandler) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/album/:album_id", h.handleRetrieveAlbum)

	withAlbumPermission := WithPermission(h.logger, h.sessionService, model.Permission_Music)

	router.HandlerFunc(http.MethodPut, "/api/album/new", withAlbumPermission(h.handleCreateAlbum))
	router.HandlerFunc(http.MethodPost, "/api/album/:album_id", withAlbumPermission(h.handleUpdateAlbum))
	router.HandlerFunc(http.MethodDelete, "/api/album/:album_id", withAlbumPermission(h.handleDeleteAlbum))
}

// //////////////////////////////////////////////////
// create muzic

func (h *albumHandler) handleCreateAlbum(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var album *model.MusicAlbum
	var err error

	switch {
	default:

		//
		// decode request
		//

		album, err = extractAlbumFromBody(req)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] create album: %#v", album))

		//
		// execute
		//

		album, err = h.service.CreateAlbum(ctx, album)
		if err != nil {
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonAlbumResponse(album))
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

func (h *albumHandler) handleRetrieveAlbum(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var albumId model.MusicAlbumId
	var album *model.MusicAlbum
	var err error

	switch {
	default:

		//
		// decode request
		//

		albumId = model.MusicAlbumId(toInt64(extractPathParameter(req, "album_id")))
		if albumId == 0 {
			err = model.ErrInvalidMusicAlbumId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] retrieve album %d", albumId))

		//
		// execute
		//

		album, err = h.service.GetAlbum(ctx, albumId)
		if err != nil {
			break
		}
		if album == nil {
			err = model.ErrMusicAlbumNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonAlbumResponse(album))
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

func (h *albumHandler) handleUpdateAlbum(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var albumId model.MusicAlbumId
	var album *model.MusicAlbum
	var err error

	switch {
	default:

		//
		// decode request
		//

		albumId = model.MusicAlbumId(toInt64(extractPathParameter(req, "album_id")))
		if albumId == 0 {
			err = model.ErrInvalidMusicAlbumId
			break
		}
		album, err = extractAlbumFromBody(req)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] update album %d: %#v", albumId, album))

		//
		// execute
		//

		album, err = h.service.UpdateAlbum(ctx, album)
		if err != nil {
			break
		}
		if album == nil {
			err = model.ErrMusicAlbumNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonAlbumResponse(album))
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

func (h *albumHandler) handleDeleteAlbum(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var albumId model.MusicAlbumId
	var err error

	switch {
	default:

		//
		// decode request
		//

		albumId = model.MusicAlbumId(toInt64(extractPathParameter(req, "album_id")))
		if albumId == 0 {
			err = model.ErrInvalidMusicAlbumId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] delete album %d", albumId))

		//
		// execute
		//

		err = h.service.DeleteAlbum(ctx, albumId)
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

func extractAlbumFromBody(req *http.Request) (*model.MusicAlbum, error) {
	var jsonBody JsonAlbumBody
	jsonErr := json.NewDecoder(req.Body).Decode(&jsonBody)
	switch {
	case jsonErr == io.EOF:
		return nil, model.ErrInvalidBody
	case jsonErr != nil:
		return nil, model.ErrInvalidBody
	}

	return toAlbum(jsonBody.Album), nil
}

type JsonAlbumBody struct {
	Album *JsonAlbumLite `json:"album,omitempty"`
}

func toAlbum(jsonAlbum *JsonAlbumLite) *model.MusicAlbum {
	return &model.MusicAlbum{
		Id:       model.MusicAlbumId(jsonAlbum.Id),
		DeezerId: model.DeezerAlbumId(jsonAlbum.DeezerId),
		Name:     jsonAlbum.Name,
		ImgUrl:   model.Url(jsonAlbum.ImgUrl),
	}
}

// //////////////////////////////////////////////////
// encode

func toJsonAlbumsResponse(albums []*model.MusicAlbum) *JsonAlbumsResponse {
	return &JsonAlbumsResponse{
		Success: true,
		Albums:  util.Convert(albums, toJsonAlbumLite),
	}
}

func toJsonAlbumLite(album *model.MusicAlbum) *JsonAlbumLite {
	if album == nil {
		return nil
	}
	return &JsonAlbumLite{
		Id:       int64(album.Id),
		DeezerId: int64(album.DeezerId),
		Name:     album.Name,
		ImgUrl:   string(album.ImgUrl),
	}
}

func toJsonAlbumResponse(album *model.MusicAlbum) *JsonAlbumResponse {
	return &JsonAlbumResponse{
		Success: true,
		Album:   toJsonAlbum(album),
	}
}

func toJsonAlbum(album *model.MusicAlbum) *JsonAlbum {
	if album == nil {
		return nil
	}
	return &JsonAlbum{
		Id:       int64(album.Id),
		DeezerId: int64(album.DeezerId),
		Name:     album.Name,
		ImgUrl:   string(album.ImgUrl),
		Musics:   util.Convert(album.Musics, toJsonMusicLite),
	}
}

type JsonAlbumsResponse struct {
	Success bool             `json:"success,omitempty"`
	Albums  []*JsonAlbumLite `json:"albums,omitempty"`
}

type JsonAlbumLite struct {
	Id       int64  `json:"id,omitempty"`
	DeezerId int64  `json:"deezerId,omitempty"`
	Name     string `json:"name,omitempty"`
	ImgUrl   string `json:"imgUrl,omitempty"`
}

type JsonAlbumResponse struct {
	Success bool       `json:"success,omitempty"`
	Album   *JsonAlbum `json:"album,omitempty"`
}

type JsonAlbum struct {
	Id       int64            `json:"id,omitempty"`
	DeezerId int64            `json:"deezerId,omitempty"`
	Name     string           `json:"name,omitempty"`
	ImgUrl   string           `json:"imgUrl,omitempty"`
	Musics   []*JsonMusicLite `json:"musics,omitempty"`
}
