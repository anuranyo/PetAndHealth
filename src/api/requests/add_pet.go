package requests

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"strings"
)

type AddPet struct {
	Name    string             `json:"name"`
	Species string             `json:"species"`
	Breed   string             `json:"breed"`
	Age     int                `json:"age"`
	OwnerID primitive.ObjectID `json:"owner_id"`
}

func NewPet(r *http.Request) (*AddPet, error) {
	bodyReader := r.Body
	if bodyReader == nil {
		return nil, errors.New("missing body")
	}

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	var pet AddPet
	err = json.Unmarshal(body, &pet)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(pet.Name) == "" {
		return nil, errors.New("pet name is required")
	}

	if strings.TrimSpace(pet.Species) == "" {
		return nil, errors.New("pet species is required")
	}

	if pet.OwnerID.IsZero() {
		return nil, errors.New("owner ID is required")
	}

	return &pet, nil
}
