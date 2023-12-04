package client_test

import (
	"testing"

	"github.com/gre-ory/amnezic-go/internal/client"
	"github.com/gre-ory/amnezic-go/internal/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSearch(t *testing.T) {

	logger := zap.L()
	deezerClient := client.NewDeezerClient(logger)
	req := model.NewDeezerSearchRequest().WithQuery("harder better").WithLimit(1)
	musics, err := deezerClient.Search(req)

	require.Equal(t, nil, err, "wrong error")
	require.Equal(t, []*model.Music{
		{
			Id:     3135556,
			Name:   "Harder, Better, Faster, Stronger",
			Mp3Url: "https://cdns-preview-d.dzcdn.net/stream/c-deda7fa9316d9e9e880d2c6207e92260-10.mp3",
			Artist: &model.MusicArtist{
				Id:     27,
				Name:   "Daft Punk",
				ImgUrl: "https://api.deezer.com/artist/27/image",
			},
			Album: &model.MusicAlbum{
				Id:     302127,
				Name:   "Discovery",
				ImgUrl: "https://api.deezer.com/album/302127/image",
			},
		},
	}, musics, "wrong musics")

}

func TestTrack(t *testing.T) {

	logger := zap.L()
	deezerClient := client.NewDeezerClient(logger)
	music, err := deezerClient.GetTrack(3135556)

	require.Equal(t, nil, err, "wrong error")
	require.Equal(t, &model.Music{
		Id:     3135556,
		Name:   "Harder, Better, Faster, Stronger",
		Mp3Url: "https://cdns-preview-d.dzcdn.net/stream/c-deda7fa9316d9e9e880d2c6207e92260-10.mp3",
		Artist: &model.MusicArtist{
			Id:     27,
			Name:   "Daft Punk",
			ImgUrl: "https://api.deezer.com/artist/27/image",
		},
		Album: &model.MusicAlbum{
			Id:     302127,
			Name:   "Discovery",
			ImgUrl: "https://api.deezer.com/album/302127/image",
		},
	}, music, "wrong music")

}
