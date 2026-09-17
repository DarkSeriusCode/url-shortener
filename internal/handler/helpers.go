package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func ReadJSONFromRequest[T Validatable](r *http.Request) (*T, error) {
	body, err := io.ReadAll(r.Body)
	if len(body) == 0 || err != nil {
		return nil, errors.New("Body cannot be empty")
	}

	var reqJSON T
	if err = json.Unmarshal(body, &reqJSON); err != nil {
		return nil, errors.New("Unable to parse json")
	}
	if err = reqJSON.Validate(); err != nil {
		return nil, errors.New(fmt.Sprintf("Validation error: %s", err))
	}
	return &reqJSON, nil
}
