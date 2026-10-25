package service

import (
	"log/slog"
	"math/big"
	"net/url"

	"github.com/DarkSeriusCode/url-shortener/internal/storage"
	"github.com/google/uuid"
	"github.com/mattheath/base62"
)

type Shortener struct {
	appStorage storage.URLStorage
}

func NewShortener(s storage.URLStorage) *Shortener {
	return &Shortener{
		appStorage: s,
	}
}

func (s *Shortener) Shorten(originURL url.URL) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		slog.Error("Error generating UUID!")
		return "", err
	}
	n := new(big.Int).SetBytes(id[:])
	urlID := base62.EncodeBigInt(n)
	err = s.appStorage.SaveUrl(originURL, urlID)
	return urlID, err
}

func (s *Shortener) Resolve(urlID string) (*storage.ShortenedURL, error) {
	return s.appStorage.GetUrl(urlID)
}
