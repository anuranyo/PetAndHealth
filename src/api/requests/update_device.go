package requests

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
)

type UpdateDevice struct {
	ID           primitive.ObjectID `json:"-"` // Тепер ID не очікується в JSON
	SerialNumber string             `json:"serial_number,omitempty"`
	Model        string             `json:"model,omitempty"`
	PetID        primitive.ObjectID `json:"pet_id,omitempty"`
	Status       string             `json:"status,omitempty"`
	LastSyncTime int64              `json:"last_sync_time,omitempty"`
}

func NewUpdateDevice(r *http.Request) (*UpdateDevice, error) {
	bodyReader := r.Body
	if bodyReader == nil {
		return nil, errors.New("missing body")
	}

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	var device UpdateDevice
	err = json.Unmarshal(body, &device)
	if err != nil {
		return nil, err
	}

	// ID тепер буде встановлено з URL параметрів, а не з тіла запиту
	return &device, nil
}
