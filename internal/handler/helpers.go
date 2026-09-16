package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadJSONFromRequest[T any](r *http.Request) (*T, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var reqJSON T
	if err = json.Unmarshal(body, &reqJSON); err != nil {
		return nil, err
	}
	return &reqJSON, nil
}
