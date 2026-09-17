package handler

import "errors"

type Validatable interface {
	Validate() error
}

// ------------------------------------------------------------------------------------------------

type URLShortenReq struct {
	URL string `json:"url" example:"https://example.com"`
}

func (r URLShortenReq) Validate() error {
	if r.URL == "" {
		return errors.New("URL is empty")
	}
	return nil
}

type URLShortenResp struct {
	ID string `json:"id"`
}

// ------------------------------------------------------------------------------------------------

type NewUserReq struct {
	Name string `json:"name"`
	Password string `json:"password"`
}

func (r NewUserReq) Validate() error {
	if r.Name == "" {
		return errors.New("User name cannot be empty")
	}
	if r.Password == "" {
		return errors.New("Password cannot be empty")
	}
	return nil
}
