package handler

type URLShortenReq struct {
	URL string `json:"url" example:"https://example.com"`
}

type URLShortenResp struct {
	ID string `json:"id"`
}
