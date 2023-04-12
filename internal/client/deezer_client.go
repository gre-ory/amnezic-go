package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// deezer client

type DeezerClient interface {
	Search(query string, limit int) ([]*model.Music, error)
	GetTrack(trackId model.DeezerMusicId) (*model.Music, error)
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
// search

func (c *deezerClient) Search(query string, limit int) ([]*model.Music, error) {

	c.logger.Info(fmt.Sprintf("[client] query: %s, limit: %d", query, limit))

	url := fmt.Sprintf("https://api.deezer.com/search?q=%s&limit=%d", url.QueryEscape(query), limit)
	c.logger.Info(fmt.Sprintf("[client] url: %s", url))

	resp, err := http.Get(url)
	if err != nil {
		c.logger.Info(fmt.Sprintf("[client] >>> error: %s", err.Error()), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	var jsonSearch JsonDeezerSearch
	err = json.NewDecoder(resp.Body).Decode(&jsonSearch)
	if err != nil {
		c.logger.Info(fmt.Sprintf("[client] >>> error: %s", err.Error()), zap.Error(err))
		return nil, err
	}

	c.logger.Info(fmt.Sprintf("[client] tracks: %d", len(jsonSearch.Tracks)))

	return util.Convert(jsonSearch.Tracks, toMusic), nil
}

// //////////////////////////////////////////////////
// track

func (c *deezerClient) GetTrack(trackId model.DeezerMusicId) (*model.Music, error) {

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

// //////////////////////////////////////////////////
// json

type JsonDeezerSearch struct {
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
