package storage

import (
	"net/url"
	"errors"
)

type StorageErr error

var (
	ErrNotFound StorageErr = errors.New("Not found")
)

// TODO: Use UUID when I back home

type Storage interface {
	SaveUrl(originURL url.URL, urlID string)
	GetUrl(urlID string) (*ShortenedURL, error)
}
