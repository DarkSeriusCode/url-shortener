package storage

import (
	"net/url"
	"time"
)

type MemoryStorage struct {
	urls map[string]ShortenedURL
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls: make(map[string]ShortenedURL),
	}
}

func (s *MemoryStorage) SaveUrl(originURL url.URL, urlID string) error {
	s.urls[urlID] = ShortenedURL{
		ID:        uint64(time.Now().UnixMilli()),
		OriginURL: &originURL,
		UrlID:     urlID,
		Author:    nil,
	}
	return nil
}

func (s *MemoryStorage) CheckExistance(originURL url.URL) (*ShortenedURL, error) {
	for _, savedURL := range s.urls {
		if savedURL.OriginURL == &originURL {
			return &savedURL, nil;
		}
	}
	return nil, nil;
}

func (s *MemoryStorage) GetUrl(urlID string) (*ShortenedURL, error) {
	entry, ok := s.urls[urlID]
	if !ok {
		return nil, ErrNotFound
	}
	return &entry, nil
}
