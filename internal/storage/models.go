package storage

import "net/url"

type ShortenedURL struct {
	ID        uint64
	OriginURL url.URL
	UrlID     string
	Author    *User
}

type User struct {
	ID       uint64
	Name     string
	Password string
}
