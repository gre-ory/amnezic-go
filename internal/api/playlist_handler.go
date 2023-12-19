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
// playlist handler

func NewPlaylisthandler(logger *zap.Logger, service service.MusicService) Handler {
	return &playlistHandler{
		logger:  logger,
		service: service,
	}
}

type playlistHandler struct {
	logger  *zap.Logger
	service service.MusicService
}

// //////////////////////////////////////////////////
// register

func (h *playlistHandler) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/deezer/playlist", h.handleSearchDeezerPlaylist)
	router.HandlerFunc(http.MethodGet, "/api/deezer/playlist/:playlist_id", h.handleRetrieveDeezerPlaylist)
}

// //////////////////////////////////////////////////
// search

func (h *playlistHandler) handleSearchDeezerPlaylist(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var search *model.SearchPlaylistRequest
	var playlists []*model.Playlist
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

		search = model.NewSearchPlaylistRequest().
			WithQuery(extractParameter(req, "search")).
			WithDeezerPlaylistId(model.DeezerPlaylistId(toInt64(extractParameter(req, "playlist_id")))).
			WithLimit(limit)

		h.logger.Info(fmt.Sprintf("[api] search %d playlists", limit), zap.Object("search", search))

		//
		// execute
		//

		playlists, err = h.service.SearchDeezerPlaylist(ctx, search)
		if err != nil {
			break
		}
		if playlists == nil {
			err = model.ErrPlaylistNotFound
			break
		}
		h.logger.Info(fmt.Sprintf("[api] found %d playlists", len(playlists)), zap.Object("search", search))

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonPlaylistsResponse(playlists))
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
// retrieve playlist

func (h *playlistHandler) handleRetrieveDeezerPlaylist(resp http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var playlistId model.DeezerPlaylistId
	var playlist *model.Playlist
	var err error

	switch {
	default:

		//
		// decode request
		//

		playlistId = model.DeezerPlaylistId(toInt64(extractPathParameter(req, "playlist_id")))
		if playlistId == 0 {
			err = model.ErrInvalidPlaylistId
			break
		}
		h.logger.Info(fmt.Sprintf("[api] retrieve playlist %d", playlistId))

		//
		// execute
		//

		playlist, err = h.service.GetDeezerPlaylist(ctx, playlistId)
		if err != nil {
			break
		}
		if playlist == nil {
			err = model.ErrPlaylistNotFound
			break
		}

		//
		// encode success
		//

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		err = json.NewEncoder(resp).Encode(toJsonPlaylistResponse(playlist))
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

func toJsonPlaylistResponse(playlist *model.Playlist) *JsonPlaylistResponse {
	return &JsonPlaylistResponse{
		Success:  true,
		Playlist: toJsonPlaylist(playlist),
	}
}

func toJsonPlaylistsResponse(playlists []*model.Playlist) *JsonPlaylistsResponse {
	return &JsonPlaylistsResponse{
		Success:   true,
		Playlists: util.Convert(playlists, toJsonPlaylistInfo),
	}
}

func toJsonPlaylistInfo(playlist *model.Playlist) *JsonPlaylistInfo {
	if playlist == nil {
		return nil
	}
	return &JsonPlaylistInfo{
		DeezerId: int64(playlist.DeezerId),
		Name:     playlist.Name,
		Public:   playlist.Public,
		ImgUrl:   playlist.ImgUrl,
		NbMusics: playlist.NbMusics,
		User:     playlist.User,
	}
}

func toJsonPlaylist(playlist *model.Playlist) *JsonPlaylist {
	if playlist == nil {
		return nil
	}
	return &JsonPlaylist{
		JsonPlaylistInfo: *toJsonPlaylistInfo(playlist),
		Musics:           util.Convert(playlist.Musics, toJsonMusic),
	}
}

type JsonPlaylistResponse struct {
	Success  bool          `json:"success,omitempty"`
	Playlist *JsonPlaylist `json:"playlist,omitempty"`
}

type JsonPlaylistsResponse struct {
	Success   bool                `json:"success,omitempty"`
	Playlists []*JsonPlaylistInfo `json:"playlists,omitempty"`
}

type JsonPlaylistInfo struct {
	DeezerId int64  `json:"deezerId,omitempty"`
	Name     string `json:"name,omitempty"`
	Public   bool   `json:"public,omitempty"`
	ImgUrl   string `json:"imgUrl,omitempty"`
	NbMusics int    `json:"nbMusics,omitempty"`
	User     string `json:"user,omitempty"`
}

type JsonPlaylist struct {
	JsonPlaylistInfo
	Musics []*JsonMusic `json:"musics,omitempty"`
}
