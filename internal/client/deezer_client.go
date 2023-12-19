package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// deezer client

type DeezerClient interface {
	SearchMusic(search *model.SearchMusicRequest) ([]*model.Music, error)
	GetMusic(trackId model.DeezerMusicId) (*model.Music, error)

	SearchPlaylist(search *model.SearchPlaylistRequest) ([]*model.Playlist, error)
	GetPlaylist(id model.DeezerPlaylistId, withTracks bool) (*model.Playlist, error)
}

func NewDeezerClient(logger *zap.Logger) DeezerClient {
	return &deezerClient{
		logger: logger,
	}
}

type deezerClient struct {
	logger *zap.Logger
}

// //////////////////////////////////////////////////
// search musics

func (c *deezerClient) SearchMusic(search *model.SearchMusicRequest) ([]*model.Music, error) {

	c.logger.Info(fmt.Sprintf("[client] search-music: req=%#v", search))

	url := fmt.Sprintf("https://api.deezer.com/search/track%s", search.ComputeParameters())
	c.logger.Info(fmt.Sprintf("[client] search-music: url=%s", url))

	resp, err := http.Get(url)
	if err != nil {
		c.logger.Info(fmt.Sprintf("[client] search-music: >>> error: %s", err.Error()), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	var jsonSearch JsonDeezerSearchTracks
	err = json.NewDecoder(resp.Body).Decode(&jsonSearch)
	if err != nil {
		c.logger.Info(fmt.Sprintf("[client] search-music: >>> error: %s", err.Error()), zap.Error(err))
		return nil, err
	}

	c.logger.Info(fmt.Sprintf("[client] search-music: %d track(s)", len(jsonSearch.Tracks)))

	return util.Convert(jsonSearch.Tracks, toMusic), nil
}

// //////////////////////////////////////////////////
// get music

func (c *deezerClient) GetMusic(trackId model.DeezerMusicId) (*model.Music, error) {

	url := fmt.Sprintf("https://api.deezer.com/track/%d", trackId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jsonTrack JsonDeezerTrack
	err = json.NewDecoder(resp.Body).Decode(&jsonTrack)
	if err != nil {
		return nil, err
	}

	return toMusic(&jsonTrack), nil
}

// //////////////////////////////////////////////////
// search playlists

func (c *deezerClient) SearchPlaylist(search *model.SearchPlaylistRequest) ([]*model.Playlist, error) {

	c.logger.Info(fmt.Sprintf("[client] search-playlist: req=%#v", search))

	playlistId := search.GetDeezerPlaylistId()
	if playlistId != 0 {
		playlist, err := c.GetPlaylist(playlistId, false /* with tracks */)
		if playlist != nil && err == nil {
			return []*model.Playlist{playlist}, nil
		}
		c.logger.Info(fmt.Sprintf("[client] search-playlist: not able to fetch playlist %d >>> fallback to normal search", playlistId), zap.Error(err))
	}

	url := fmt.Sprintf("https://api.deezer.com/search/playlist%s", search.ComputeParameters())
	c.logger.Info(fmt.Sprintf("[client] search-playlist: url=%s", url))

	resp, err := http.Get(url)
	if err != nil {
		c.logger.Info(fmt.Sprintf("[client] search-playlist: >>> error: %s", err.Error()), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	var jsonResponse JsonDeezerSearchPlaylists
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		c.logger.Info(fmt.Sprintf("[client] search-playlist: >>> error: %s", err.Error()), zap.Error(err))
		return nil, err
	}

	c.logger.Info(fmt.Sprintf("[client] search-playlist: %d playlist(s)", len(jsonResponse.Playlists)))

	return util.Convert(jsonResponse.Playlists, toPlaylist), nil
}

// //////////////////////////////////////////////////
// get playlist

func (c *deezerClient) GetPlaylist(id model.DeezerPlaylistId, withTracks bool) (*model.Playlist, error) {

	c.logger.Info(fmt.Sprintf("[client] get-playlist: id=%d, withTracks=%t", id, withTracks))

	url := fmt.Sprintf("https://api.deezer.com/playlist/%d", id)
	c.logger.Info(fmt.Sprintf("[client] get-playlist: url=%s", url))

	resp, err := http.Get(url)
	if err != nil {
		c.logger.Info(fmt.Sprintf("[client] get-playlist: >>> error: %s", err.Error()), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	var jsonResponse JsonDeezerPlaylist
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		c.logger.Info(fmt.Sprintf("[client] get-playlist: >>> error: %s", err.Error()), zap.Error(err))
		return nil, err
	}

	c.logger.Info(fmt.Sprintf("[client] get-playlist: title=%s", jsonResponse.Title))

	playlist := toPlaylist(&jsonResponse)
	if withTracks && jsonResponse.Tracks != nil {
		playlist.Musics = util.Convert(jsonResponse.Tracks.Tracks, toMusic)
	}
	return playlist, nil
}

// //////////////////////////////////////////////////
// adapter

func toMusic(jsonTrack *JsonDeezerTrack) *model.Music {
	return &model.Music{
		DeezerId: model.DeezerMusicId(jsonTrack.Id),
		Name:     jsonTrack.Title,
		Mp3Url:   jsonTrack.Preview,
		Artist:   toArtist(jsonTrack.Artist),
		Album:    toAlbum(jsonTrack.Album),
	}
}

func toArtist(jsonArtist *JsonDeezerArtist) *model.MusicArtist {
	if jsonArtist == nil {
		return nil
	}
	return &model.MusicArtist{
		DeezerId: model.DeezerArtistId(jsonArtist.Id),
		Name:     jsonArtist.Name,
		ImgUrl:   jsonArtist.Picture,
	}
}

func toAlbum(jsonAlbum *JsonDeezerAlbum) *model.MusicAlbum {
	if jsonAlbum == nil {
		return nil
	}
	return &model.MusicAlbum{
		DeezerId: model.DeezerAlbumId(jsonAlbum.Id),
		Name:     jsonAlbum.Title,
		ImgUrl:   jsonAlbum.Cover,
	}
}

func toPlaylist(jsonPlaylist *JsonDeezerPlaylist) *model.Playlist {
	playlist := &model.Playlist{
		DeezerId:    model.DeezerPlaylistId(jsonPlaylist.Id),
		Name:        jsonPlaylist.Title,
		Public:      jsonPlaylist.Public,
		PlaylistUrl: jsonPlaylist.Link,
		ImgUrl:      jsonPlaylist.Picture,
		NbMusics:    jsonPlaylist.NbTracks,
	}
	if jsonPlaylist.Creator != nil {
		playlist.User = jsonPlaylist.Creator.Name
	} else if jsonPlaylist.User != nil {
		playlist.User = jsonPlaylist.User.Name
	}
	return playlist
}

// //////////////////////////////////////////////////
// json

type JsonDeezerSearchTracks struct {
	Tracks []*JsonDeezerTrack `json:"data"`
}

type JsonDeezerTrack struct {
	Id                 int64                    `json:"id"`
	Readable           bool                     `json:"readable"`
	Title              string                   `json:"title"`
	TitleShort         string                   `json:"title_short"`
	Duration           int64                    `json:"duration"`
	Preview            string                   `json:"preview"`
	AvailableCountries []string                 `json:"available_countries"`
	Contributors       []*JsonDeezerContributor `json:"contributors"`
	Artist             *JsonDeezerArtist        `json:"artist"`
	Album              *JsonDeezerAlbum         `json:"album"`
}

type JsonDeezerContributor struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Type    string `json:"type"`
	Role    string `json:"role"`
}

type JsonDeezerArtist struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Type    string `json:"type"`
	Role    string `json:"role"`
}

type JsonDeezerAlbum struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Type  string `json:"type"`
	Role  string `json:"role"`
}

type JsonDeezerSearchPlaylists struct {
	Playlists []*JsonDeezerPlaylist `json:"data"`
}

type JsonDeezerPlaylist struct {
	Id        int64                     `json:"id"`
	Title     string                    `json:"title"`
	Public    bool                      `json:"public"`
	NbTracks  int                       `json:"nb_tracks"`
	Link      string                    `json:"link"`
	Picture   string                    `json:"picture"`
	TrackList string                    `json:"tracklist"`
	Creator   *JsonDeezerCreator        `json:"creator"`
	User      *JsonDeezerUser           `json:"user"`
	Tracks    *JsonDeezerPlaylistTracks `json:"tracks"`
}

type JsonDeezerCreator struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type JsonDeezerUser struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type JsonDeezerPlaylistTracks struct {
	Tracks []*JsonDeezerTrack `json:"data"`
}
