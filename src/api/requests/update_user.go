package requests

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
)

type UpdateUser struct {
	ID       primitive.ObjectID   `json:"-"` // ID тепер не очікується в JSON
	FullName string               `json:"full_name,omitempty"`
	Role     string               `json:"role,omitempty"`
	Email    string               `json:"email,omitempty"`
	Password string               `json:"password,omitempty"`
	PetsID   []primitive.ObjectID `json:"pets_id,omitempty"`
}

func NewUpdateUser(r *http.Request) (*UpdateUser, error) {
	bodyReader := r.Body
	if bodyReader == nil {
		return nil, errors.New("missing body")
	}

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	var user UpdateUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
