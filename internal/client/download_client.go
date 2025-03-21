package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gre-ory/amnezic-go/internal/model"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// download client

type DownloadClient interface {
	DownloadMusic(url model.Url, fileName model.Url) (err error)
	DownloadImage(url model.Url, fileName model.Url) (err error)
}

func NewDownloadClient(logger *zap.Logger, musicFilter *model.FileFilter, imageFilter *model.FileFilter) DownloadClient {
	return &downloadClient{
		logger:      logger,
		musicFilter: musicFilter,
		imageFilter: imageFilter,
	}
}

type downloadClient struct {
	logger      *zap.Logger
	musicFilter *model.FileFilter
	imageFilter *model.FileFilter
}

// //////////////////////////////////////////////////
// download

func (c *downloadClient) DownloadMusic(url model.Url, fileName model.Url) error {
	return c.downloadFile(url, fileName, c.musicFilter)
}

func (c *downloadClient) DownloadImage(url model.Url, fileName model.Url) error {
	return c.downloadFile(url, fileName, c.imageFilter)
}

func (c *downloadClient) downloadFile(url model.Url, fileName model.Url, filter *model.FileFilter) error {

	filePath := filepath.Join(filter.Directory, string(fileName))
	if !filter.MatchExtension(filePath) {
		return model.ErrInvalidExtension
	}

	if _, err := os.Stat(filePath); err == nil {
		return model.ErrFileAlreadyExists(string(fileName))
	}

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	return c.download(url, out)
}

func (c *downloadClient) download(url model.Url, writer io.Writer) (err error) {

	// Get the data

	resp, err := http.Get(string(url))
	if err != nil {
		c.logger.Info(fmt.Sprintf("[ KO ] download from url %s", url), zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
		c.logger.Info(fmt.Sprintf("[ KO ] download from url %s", url), zap.Error(err))
		return err
	}

	// Writer the body to file
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		c.logger.Info(fmt.Sprintf("[ KO ] download from url %s", url), zap.Error(err))
		return err
	}

	c.logger.Info(fmt.Sprintf("[ OK ] download from url %s", url))
	return nil
}
