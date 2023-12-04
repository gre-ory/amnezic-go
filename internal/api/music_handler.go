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
	router.HandlerFunc(http.MethodPost, "/api/music/:music_id", h.handleUpdateMusic)
	router.HandlerFunc(http.MethodDelete, "/api/music/:music_id", h.handleDeleteMusic)
}

// //////////////////////////////////////////////////
// search

func (h *musicHandler) handleSearchMusic(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var search *model.DeezerSearch
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

		search = model.NewDeezerSearchRequest().
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

		musics, err = h.service.SearchMusic(ctx, search)
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

func toMusic(jsonMusic *JsonMusic) *model.Music {

	artist := toArtist(jsonMusic.Artist)
	album := toAlbum(jsonMusic.Album)

	music := &model.Music{
		Id:       model.MusicId(jsonMusic.Id),
		DeezerId: model.DeezerMusicId(jsonMusic.DeezerId),
		Name:     jsonMusic.Name,
		Mp3Url:   jsonMusic.Mp3Url,
		ArtistId: artist.Id,
		Artist:   artist,
		AlbumId:  album.Id,
		Album:    album,
	}

	return music
}

func toArtist(jsonArtist *JsonMusicArtist) *model.MusicArtist {
	return &model.MusicArtist{
		Id:       model.MusicArtistId(jsonArtist.Id),
		DeezerId: model.DeezerArtistId(jsonArtist.DeezerId),
		Name:     jsonArtist.Name,
		ImgUrl:   jsonArtist.ImgUrl,
	}
}

func toAlbum(jsonAlbum *JsonMusicAlbum) *model.MusicAlbum {
	return &model.MusicAlbum{
		Id:       model.MusicAlbumId(jsonAlbum.Id),
		DeezerId: model.DeezerAlbumId(jsonAlbum.DeezerId),
		Name:     jsonAlbum.Name,
		ImgUrl:   jsonAlbum.ImgUrl,
	}
}

type JsonMusicBody struct {
	Music *JsonMusic `json:"music,omitempty"`
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
		Id:        int64(music.Id),
		DeezerId:  int64(music.DeezerId),
		Name:      music.Name,
		Mp3Url:    music.Mp3Url,
		Artist:    toJsonArtist(music.Artist),
		Album:     toJsonAlbum(music.Album),
		Questions: util.Convert(music.Questions, toJsonThemeQuestion),
	}
}

func toJsonArtist(artist *model.MusicArtist) *JsonMusicArtist {
	if artist == nil {
		return nil
	}
	return &JsonMusicArtist{
		Id:       int64(artist.Id),
		DeezerId: int64(artist.DeezerId),
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
		DeezerId: int64(album.DeezerId),
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
	Id        int64                `json:"id,omitempty"`
	DeezerId  int64                `json:"deezerId,omitempty"`
	Name      string               `json:"name,omitempty"`
	Mp3Url    string               `json:"mp3Url,omitempty"`
	Artist    *JsonMusicArtist     `json:"artist,omitempty"`
	Album     *JsonMusicAlbum      `json:"album,omitempty"`
	Questions []*JsonThemeQuestion `json:"questions,omitempty"`
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
