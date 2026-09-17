package service

import (
	"crypto/rand"
	"log/slog"
	"math/big"
	"net/url"
	"time"

	"github.com/DarkSeriusCode/url-shortener/internal/storage"
	"github.com/mattheath/base62"
)

type Shortener struct {
	appStorage storage.Storage
}

func NewShortener(s storage.Storage) *Shortener {
	return &Shortener{
		appStorage: s,
	}
}

func (s *Shortener) Shorten(originURL url.URL) (string, error) {
	// TODO: Use UUID!!!!!
	timestamp := big.NewInt(time.Now().UnixMilli())
	noice, err := rand.Int(rand.Reader, timestamp)
	if err != nil {
		slog.Error("Unexpected error with noice generation")
	}
	urlID := base62.EncodeBigInt(timestamp.Add(timestamp, noice))
	err = s.appStorage.SaveUrl(originURL, urlID)
	return urlID, err
}

func (s *Shortener) Resolve(urlID string) (*storage.ShortenedURL, error) {
	return s.appStorage.GetUrl(urlID)
}
