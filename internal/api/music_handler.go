package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/service"
	"github.com/gre-ory/amnezic-go/internal/util"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// music handler

func NewMusichandler(logger *zap.Logger, service service.MusicService) Handler {
	return &musicHandler{
		logger:  logger,
		service: service,
	}
}

type musicHandler struct {
	logger  *zap.Logger
	service service.MusicService
}

// //////////////////////////////////////////////////
// register

func (h *musicHandler) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/music", h.handleSearchMusic)
	router.HandlerFunc(http.MethodPut, "/api/music/new", h.handleCreateMusic)
	router.HandlerFunc(http.MethodGet, "/api/music/:music_id", h.handleRetrieveMusic)
	router.HandlerFunc(http.MethodDelete, "/api/music/:music_id", h.handleDeleteMusic)
}

// //////////////////////////////////////////////////
// search

func (h *musicHandler) handleSearchMusic(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var query string
	var limit int
	var musics []*model.Music
	var err error

	switch {
	default:

		//
		// decode request
		//

		query = extractParameter(req, "search")
		limit = toInt(extractParameter(req, "limit"))
		if limit == 0 {
			limit = 42
		}
		h.logger.Info(fmt.Sprintf("[api] search %d musics matching %s", limit, query))

		//
		// execute
		//

		musics, err = h.service.SearchMusic(ctx, query, limit)
		if err != nil {
			break
		}
		if musics == nil {
			err = model.ErrMusicNotFound
			break
		}
		h.logger.Info(fmt.Sprintf("[api] found %d musics matching %s", len(musics), query))

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonMusicsResponse(musics))
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

func (h *musicHandler) handleCreateMusic(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var deezerId model.DeezerMusicId
	var music *model.Music
	var err error

	switch {
	default:

		//
		// decode request
		//

		deezerId = model.DeezerMusicId(toInt64(extractParameter(req, "deezer_id")))
		if deezerId == 0 {
			err = model.ErrInvalidDeezerId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] create music from deezer-id %d", deezerId))

		//
		// execute
		//

		music, err = h.service.AddDeezerMusic(ctx, deezerId)
		if err != nil {
			break
		}
		if music == nil {
			err = model.ErrMusicNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonMusicResponse(music))
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

func (h *musicHandler) handleRetrieveMusic(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var musicId model.MusicId
	var music *model.Music
	var err error

	switch {
	default:

		//
		// decode request
		//

		musicId = model.MusicId(toInt64(extractPathParameter(req, "music_id")))
		if musicId == 0 {
			err = model.ErrInvalidMusicId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] retrieve music %d", musicId))

		//
		// execute
		//

		music, err = h.service.GetMusic(ctx, musicId)
		if err != nil {
			break
		}
		if music == nil {
			err = model.ErrMusicNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonMusicResponse(music))
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

func (h *musicHandler) handleDeleteMusic(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var musicId model.MusicId
	var err error

	switch {
	default:

		//
		// decode request
		//

		musicId = model.MusicId(toInt64(extractPathParameter(req, "music_id")))
		if musicId == 0 {
			err = model.ErrInvalidMusicId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] delete music %d", musicId))

		//
		// execute
		//

		err = h.service.DeleteMusic(ctx, musicId)
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

func toJsonMusicResponse(music *model.Music) *JsonMusicResponse {
	return &JsonMusicResponse{
		Success: true,
		Music:   toJsonMusic(music),
	}
}

func toJsonMusicsResponse(musics []*model.Music) *JsonMusicsResponse {
	return &JsonMusicsResponse{
		Success: true,
		Musics:  util.Convert(musics, toJsonMusic),
	}
}

func toJsonMusic(music *model.Music) *JsonMusic {
	if music == nil {
		return nil
	}
	return &JsonMusic{
		Id:       int64(music.Id),
		DeezerId: int64(music.DeezerMusicId),
		Name:     music.Name,
		Mp3Url:   music.Mp3Url,
		Artist:   toJsonArtist(music.Artist),
		Album:    toJsonAlbum(music.Album),
	}
}

func toJsonArtist(artist *model.MusicArtist) *JsonMusicArtist {
	if artist == nil {
		return nil
	}
	return &JsonMusicArtist{
		Id:       int64(artist.Id),
		DeezerId: int64(artist.DeezerArtistId),
		Name:     artist.Name,
		ImgUrl:   artist.ImgUrl,
	}
}

func toJsonAlbum(album *model.MusicAlbum) *JsonMusicAlbum {
	if album == nil {
		return nil
	}
	return &JsonMusicAlbum{
		Id:       int64(album.Id),
		DeezerId: int64(album.DeezerAlbumId),
		Name:     album.Name,
		ImgUrl:   album.ImgUrl,
	}
}

type JsonMusicResponse struct {
	Success bool       `json:"success,omitempty"`
	Music   *JsonMusic `json:"music,omitempty"`
}

type JsonMusicsResponse struct {
	Success bool         `json:"success,omitempty"`
	Musics  []*JsonMusic `json:"musics,omitempty"`
}

type JsonMusic struct {
	Id       int64            `json:"id,omitempty"`
	DeezerId int64            `json:"deezerId,omitempty"`
	Name     string           `json:"name,omitempty"`
	Mp3Url   string           `json:"mp3Url,omitempty"`
	Artist   *JsonMusicArtist `json:"artist,omitempty"`
	Album    *JsonMusicAlbum  `json:"album,omitempty"`
}

type JsonMusicArtist struct {
	Id       int64  `json:"id,omitempty"`
	DeezerId int64  `json:"deezerId,omitempty"`
	Name     string `json:"name,omitempty"`
	ImgUrl   string `json:"imgUrl,omitempty"`
}

type JsonMusicAlbum struct {
	Id       int64  `json:"id,omitempty"`
	DeezerId int64  `json:"deezerId,omitempty"`
	Name     string `json:"name,omitempty"`
	ImgUrl   string `json:"imgUrl,omitempty"`
}