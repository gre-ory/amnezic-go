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
// artist handler

func NewArtisthandler(logger *zap.Logger, service service.ArtistService, sessionService service.SessionService) Handler {
	return &artistHandler{
		logger:         logger,
		service:        service,
		sessionService: sessionService,
	}
}

type artistHandler struct {
	logger         *zap.Logger
	service        service.ArtistService
	sessionService service.SessionService
}

// //////////////////////////////////////////////////
// register

func (h *artistHandler) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/artist/:artist_id", h.handleRetrieveArtist)

	withArtistPermission := WithPermission(h.logger, h.sessionService, model.Permission_Music)

	router.HandlerFunc(http.MethodPut, "/api/artist/new", withArtistPermission(h.handleCreateArtist))
	router.HandlerFunc(http.MethodPost, "/api/artist/:artist_id", withArtistPermission(h.handleUpdateArtist))
	router.HandlerFunc(http.MethodDelete, "/api/artist/:artist_id", withArtistPermission(h.handleDeleteArtist))
}

// //////////////////////////////////////////////////
// create muzic

func (h *artistHandler) handleCreateArtist(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var artist *model.MusicArtist
	var err error

	switch {
	default:

		//
		// decode request
		//

		artist, err = extractArtistFromBody(req)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] create artist: %#v", artist))

		//
		// execute
		//

		artist, err = h.service.CreateArtist(ctx, artist)
		if err != nil {
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonArtistResponse(artist))
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

func (h *artistHandler) handleRetrieveArtist(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var artistId model.MusicArtistId
	var artist *model.MusicArtist
	var err error

	switch {
	default:

		//
		// decode request
		//

		artistId = model.MusicArtistId(toInt64(extractPathParameter(req, "artist_id")))
		if artistId == 0 {
			err = model.ErrInvalidMusicArtistId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] retrieve artist %d", artistId))

		//
		// execute
		//

		artist, err = h.service.GetArtist(ctx, artistId)
		if err != nil {
			break
		}
		if artist == nil {
			err = model.ErrMusicArtistNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonArtistResponse(artist))
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

func (h *artistHandler) handleUpdateArtist(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var artistId model.MusicArtistId
	var artist *model.MusicArtist
	var err error

	switch {
	default:

		//
		// decode request
		//

		artistId = model.MusicArtistId(toInt64(extractPathParameter(req, "artist_id")))
		if artistId == 0 {
			err = model.ErrInvalidMusicArtistId
			break
		}
		artist, err = extractArtistFromBody(req)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] update artist %d: %#v", artistId, artist))

		//
		// execute
		//

		artist, err = h.service.UpdateArtist(ctx, artist)
		if err != nil {
			break
		}
		if artist == nil {
			err = model.ErrMusicArtistNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonArtistResponse(artist))
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

func (h *artistHandler) handleDeleteArtist(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var artistId model.MusicArtistId
	var err error

	switch {
	default:

		//
		// decode request
		//

		artistId = model.MusicArtistId(toInt64(extractPathParameter(req, "artist_id")))
		if artistId == 0 {
			err = model.ErrInvalidMusicArtistId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] delete artist %d", artistId))

		//
		// execute
		//

		err = h.service.DeleteArtist(ctx, artistId)
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

func extractArtistFromBody(req *http.Request) (*model.MusicArtist, error) {
	var jsonBody JsonArtistBody
	jsonErr := json.NewDecoder(req.Body).Decode(&jsonBody)
	switch {
	case jsonErr == io.EOF:
		return nil, model.ErrInvalidBody
	case jsonErr != nil:
		return nil, model.ErrInvalidBody
	}

	return toArtist(jsonBody.Artist), nil
}

type JsonArtistBody struct {
	Artist *JsonArtistLite `json:"artist,omitempty"`
}

func toArtist(jsonArtist *JsonArtistLite) *model.MusicArtist {
	return &model.MusicArtist{
		Id:       model.MusicArtistId(jsonArtist.Id),
		DeezerId: model.DeezerArtistId(jsonArtist.DeezerId),
		Name:     jsonArtist.Name,
		ImgUrl:   model.Url(jsonArtist.ImgUrl),
	}
}

// //////////////////////////////////////////////////
// encode

func toJsonArtistsResponse(artists []*model.MusicArtist) *JsonArtistsResponse {
	return &JsonArtistsResponse{
		Success: true,
		Artists: util.Convert(artists, toJsonArtistLite),
	}
}

func toJsonArtistLite(artist *model.MusicArtist) *JsonArtistLite {
	if artist == nil {
		return nil
	}
	return &JsonArtistLite{
		Id:       int64(artist.Id),
		DeezerId: int64(artist.DeezerId),
		Name:     artist.Name,
		ImgUrl:   string(artist.ImgUrl),
	}
}

func toJsonArtistResponse(artist *model.MusicArtist) *JsonArtistResponse {
	return &JsonArtistResponse{
		Success: true,
		Artist:  toJsonArtist(artist),
	}
}

func toJsonArtist(artist *model.MusicArtist) *JsonArtist {
	if artist == nil {
		return nil
	}
	return &JsonArtist{
		Id:       int64(artist.Id),
		DeezerId: int64(artist.DeezerId),
		Name:     artist.Name,
		ImgUrl:   string(artist.ImgUrl),
		Musics:   util.Convert(artist.Musics, toJsonMusicLite),
	}
}

type JsonArtistsResponse struct {
	Success bool              `json:"success,omitempty"`
	Artists []*JsonArtistLite `json:"artists,omitempty"`
}

type JsonArtistLite struct {
	Id       int64  `json:"id,omitempty"`
	DeezerId int64  `json:"deezerId,omitempty"`
	Name     string `json:"name,omitempty"`
	ImgUrl   string `json:"imgUrl,omitempty"`
}

type JsonArtistResponse struct {
	Success bool        `json:"success,omitempty"`
	Artist  *JsonArtist `json:"artist,omitempty"`
}

type JsonArtist struct {
	Id       int64            `json:"id,omitempty"`
	DeezerId int64            `json:"deezerId,omitempty"`
	Name     string           `json:"name,omitempty"`
	ImgUrl   string           `json:"imgUrl,omitempty"`
	Musics   []*JsonMusicLite `json:"musics,omitempty"`
}
