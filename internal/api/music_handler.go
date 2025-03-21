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
// music handler

func NewMusichandler(logger *zap.Logger, service service.MusicService, sessionService service.SessionService) Handler {
	return &musicHandler{
		logger:         logger,
		service:        service,
		sessionService: sessionService,
	}
}

type musicHandler struct {
	logger         *zap.Logger
	service        service.MusicService
	sessionService service.SessionService
}

// //////////////////////////////////////////////////
// register

func (h *musicHandler) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/deezer/music", h.handleSearchDeezerMusic)
	router.HandlerFunc(http.MethodGet, "/api/music/:music_id", h.handleRetrieveMusic)

	withMusicPermission := WithPermission(h.logger, h.sessionService, model.Permission_Music)

	router.HandlerFunc(http.MethodPut, "/api/music/new", withMusicPermission(h.handleCreateMusic))
	router.HandlerFunc(http.MethodPost, "/api/music/:music_id", withMusicPermission(h.handleUpdateMusic))
	router.HandlerFunc(http.MethodDelete, "/api/music/:music_id", withMusicPermission(h.handleDeleteMusic))
}

// //////////////////////////////////////////////////
// search

func (h *musicHandler) handleSearchDeezerMusic(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var search *model.SearchDeezerMusicRequest
	var musics []*model.Music
	var err error

	switch {
	default:

		//
		// decode request
		//

		limit := toInt(extractParameter(req, "limit"))
		if limit == 0 {
			limit = 100
		}

		search = model.NewSearchDeezerMusicRequest().
			WithQuery(extractParameter(req, "search")).
			WithAlbum(extractParameter(req, "album")).
			WithArtist(extractParameter(req, "artist")).
			WithLabel(extractParameter(req, "label")).
			WithTrack(extractParameter(req, "artist")).
			WithLimit(limit)

		h.logger.Info(fmt.Sprintf("[api] search %d musics", limit), zap.Object("search", search))

		//
		// execute
		//

		musics, err = h.service.SearchDeezerMusic(ctx, search)
		if err != nil {
			break
		}
		if musics == nil {
			err = model.ErrMusicNotFound
			break
		}
		h.logger.Info(fmt.Sprintf("[api] found %d musics", len(musics)), zap.Object("search", search))

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
// add deezer muzic

func (h *musicHandler) handleAddDezzerMusic(resp http.ResponseWriter, req *http.Request) {

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
		h.logger.Info(fmt.Sprintf("[api] add deezer music %d", deezerId))

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
// create muzic

func (h *musicHandler) handleCreateMusic(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var music *model.Music
	var err error

	switch {
	default:

		//
		// decode request
		//

		music, err = extractMusicFromBody(req)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] create music: %#v", music))

		//
		// execute
		//

		music, err = h.service.CreateMusic(ctx, music)
		if err != nil {
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
// update

func (h *musicHandler) handleUpdateMusic(resp http.ResponseWriter, req *http.Request) {

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
		music, err = extractMusicFromBody(req)
		if err != nil {
			break
		}
		h.logger.Info(fmt.Sprintf("[api] update music %d: %#v", musicId, music))

		//
		// execute
		//

		music, err = h.service.UpdateMusic(ctx, music)
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
// decode

func extractMusicFromBody(req *http.Request) (*model.Music, error) {
	var jsonBody JsonMusicBody
	jsonErr := json.NewDecoder(req.Body).Decode(&jsonBody)
	switch {
	case jsonErr == io.EOF:
		return nil, model.ErrInvalidBody
	case jsonErr != nil:
		return nil, model.ErrInvalidBody
	}

	return toMusic(jsonBody.Music), nil
}

type JsonMusicBody struct {
	Music *JsonMusic `json:"music,omitempty"`
}

func toMusic(jsonMusic *JsonMusic) *model.Music {

	artist := toArtist(jsonMusic.Artist)
	album := toAlbum(jsonMusic.Album)

	music := &model.Music{
		Id:       model.MusicId(jsonMusic.Id),
		DeezerId: model.DeezerMusicId(jsonMusic.DeezerId),
		Name:     jsonMusic.Name,
		Mp3Url:   model.Url(jsonMusic.Mp3Url),
		ArtistId: artist.Id,
		AlbumId:  album.Id,

		Artist: artist,
		Album:  album,
	}

	return music
}

// //////////////////////////////////////////////////
// encode

func toJsonMusicsResponse(musics []*model.Music) *JsonMusicsResponse {
	return &JsonMusicsResponse{
		Success: true,
		Musics:  util.Convert(musics, toJsonMusicLite),
	}
}

func toJsonMusicLite(music *model.Music) *JsonMusicLite {
	if music == nil {
		return nil
	}
	return &JsonMusicLite{
		Id:       int64(music.Id),
		DeezerId: int64(music.DeezerId),
		Name:     music.Name,
		Mp3Url:   string(music.Mp3Url),
		ArtistId: int64(music.ArtistId),
		AlbumId:  int64(music.AlbumId),
	}
}

func toJsonMusicResponse(music *model.Music) *JsonMusicResponse {
	return &JsonMusicResponse{
		Success: true,
		Music:   toJsonMusic(music),
	}
}

func toJsonMusic(music *model.Music) *JsonMusic {
	if music == nil {
		return nil
	}
	return &JsonMusic{
		Id:        int64(music.Id),
		DeezerId:  int64(music.DeezerId),
		Name:      music.Name,
		Mp3Url:    string(music.Mp3Url),
		Artist:    toJsonArtistLite(music.Artist),
		Album:     toJsonAlbumLite(music.Album),
		Questions: util.Convert(music.Questions, toJsonThemeQuestion),
	}
}

type JsonMusicsResponse struct {
	Success bool             `json:"success,omitempty"`
	Musics  []*JsonMusicLite `json:"musics,omitempty"`
}

type JsonMusicLite struct {
	Id       int64  `json:"id,omitempty"`
	DeezerId int64  `json:"deezerId,omitempty"`
	Name     string `json:"name,omitempty"`
	Mp3Url   string `json:"mp3Url,omitempty"`
	ArtistId int64  `json:"artistId,omitempty"`
	AlbumId  int64  `json:"albumId,omitempty"`
}

type JsonMusicResponse struct {
	Success bool       `json:"success,omitempty"`
	Music   *JsonMusic `json:"music,omitempty"`
}

type JsonMusic struct {
	Id        int64                `json:"id,omitempty"`
	DeezerId  int64                `json:"deezerId,omitempty"`
	Name      string               `json:"name,omitempty"`
	Mp3Url    string               `json:"mp3Url,omitempty"`
	Artist    *JsonArtistLite      `json:"artist,omitempty"`
	Album     *JsonAlbumLite       `json:"album,omitempty"`
	Questions []*JsonThemeQuestion `json:"questions,omitempty"`
}
