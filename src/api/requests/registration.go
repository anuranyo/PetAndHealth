package requests

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Registration struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewRegistration(r *http.Request) (*Registration, error) {
	bodyReader := r.Body
	if bodyReader == nil {
		return nil, errors.New("missing body")
	}

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	var user Registration
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
