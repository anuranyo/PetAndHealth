package requests

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type UpdateUserProfile struct {
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func NewUpdateUserProfile(r *http.Request) (*UpdateUserProfile, error) {
	bodyReader := r.Body
	if bodyReader == nil {
		return nil, errors.New("missing body")
	}

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	var profile UpdateUserProfile
	err = json.Unmarshal(body, &profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
