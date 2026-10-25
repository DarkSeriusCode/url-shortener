package storage

import (
	"errors"
	"net/url"
)

type StorageErr error

var (
	ErrNotFound StorageErr = errors.New("Not found")
)

type URLStorage interface {
	SaveUrl(originURL url.URL, urlID string) error
	CheckExistance(originURL url.URL) (*ShortenedURL, error)
	GetUrl(urlID string) (*ShortenedURL, error)
}
