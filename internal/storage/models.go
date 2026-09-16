package storage

import "net/url"

type ShortenedURL struct {
	ID uint64
	originURL url.URL
	urlID string
	author *User
}

type User struct {
	ID uint64
	name string
	password string
}
