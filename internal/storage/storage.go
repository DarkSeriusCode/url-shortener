package storage

import (
	"errors"
	"net/url"
)

type StorageErr error

var (
	ErrNotFound StorageErr = errors.New("Not found")
)

// TODO: Use UUID when I back home

type Storage interface {
	NewUser(name string, password string) error
	SaveUrl(originURL url.URL, urlID string) error
	GetUrl(urlID string) (*ShortenedURL, error)
}
