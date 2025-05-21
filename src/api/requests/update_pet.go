package requests

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
)

type UpdatePet struct {
	ID      primitive.ObjectID `json:"-"`
	Name    string             `json:"name"`
	Species string             `json:"species"`
	Breed   string             `json:"breed"`
	Age     int                `json:"age"`
	OwnerID primitive.ObjectID `json:"owner_id"`
}

func NewUpdatePet(r *http.Request) (*UpdatePet, error) {
	bodyReader := r.Body
	if bodyReader == nil {
		return nil, errors.New("missing body")
	}

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	var pet UpdatePet
	err = json.Unmarshal(body, &pet)
	if err != nil {
		return nil, err
	}

	return &pet, nil
}
