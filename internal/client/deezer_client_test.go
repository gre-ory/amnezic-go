package client_test

import (
	"testing"

	"github.com/gre-ory/amnezic-go/internal/client"
	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSearchMusic(t *testing.T) {

	logger := zap.L()
	deezerClient := client.NewDeezerClient(logger)
	req := model.NewSearchDeezerMusicRequest().WithQuery("harder better").WithLimit(1).WithStrict(true)
	musics, err := deezerClient.SearchMusic(req)

	require.Equal(t, nil, err, "wrong error")
	require.Equal(t, []*model.Music{
		{
			DeezerId: 3135556,
			Name:     "Harder, Better, Faster, Stronger",
			Mp3Url:   model.Url("https://cdns-preview-d.dzcdn.net/stream/c-deda7fa9316d9e9e880d2c6207e92260-10.mp3"),
			Artist: &model.MusicArtist{
				DeezerId: 27,
				Name:     "Daft Punk",
				ImgUrl:   model.Url("https://api.deezer.com/artist/27/image"),
			},
			Album: &model.MusicAlbum{
				DeezerId: 302127,
				Name:     "Discovery",
				ImgUrl:   model.Url("https://api.deezer.com/album/302127/image"),
			},
		},
	}, musics, "wrong musics")

}

func TestGetMusic(t *testing.T) {

	logger := zap.L()
	deezerClient := client.NewDeezerClient(logger)
	music, err := deezerClient.GetMusic(3135556)

	require.Equal(t, nil, err, "wrong error")
	require.Equal(t, &model.Music{
		DeezerId: 3135556,
		Name:     "Harder, Better, Faster, Stronger",
		Mp3Url:   model.Url("https://cdns-preview-d.dzcdn.net/stream/c-deda7fa9316d9e9e880d2c6207e92260-10.mp3"),
		Artist: &model.MusicArtist{
			DeezerId: 27,
			Name:     "Daft Punk",
			ImgUrl:   model.Url("https://api.deezer.com/artist/27/image"),
		},
		Album: &model.MusicAlbum{
			DeezerId: 302127,
			Name:     "Discovery",
			ImgUrl:   model.Url("https://api.deezer.com/album/302127/image"),
		},
	}, music, "wrong music")

}

func TestSearchPlaylist(t *testing.T) {

	logger := zap.L()
	deezerClient := client.NewDeezerClient(logger)
	req := model.NewSearchDeezerPlaylistRequest().WithQuery("eminem").WithLimit(1).WithStrict(true)
	playlists, err := deezerClient.SearchPlaylist(req)

	require.Equal(t, nil, err, "wrong error")
	require.Equal(t, []*model.Playlist{
		{
			DeezerId:    3645740262,
			Name:        "100% Eminem",
			Public:      true,
			PlaylistUrl: "https://www.deezer.com/playlist/3645740262",
			ImgUrl:      "https://api.deezer.com/playlist/3645740262/image",
			NbMusics:    39,
			User:        "Deezer Artist Editor",
		},
	}, playlists, "wrong playlists")

}

func TestGetPlaylist(t *testing.T) {

	logger := zap.L()
	deezerClient := client.NewDeezerClient(logger)
	playlist, err := deezerClient.GetPlaylist(3645740262, false)

	require.Equal(t, nil, err, "wrong error")
	require.Equal(t, &model.Playlist{
		DeezerId:    3645740262,
		Name:        "100% Eminem",
		Public:      true,
		PlaylistUrl: "https://www.deezer.com/playlist/3645740262",
		ImgUrl:      "https://api.deezer.com/playlist/3645740262/image",
		NbMusics:    39,
		User:        "Deezer Artist Editor",
	}, playlist, "wrong playlist")

}
