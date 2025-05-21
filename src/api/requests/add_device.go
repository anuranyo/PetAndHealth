package requests

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"strings"
	"time"
)

type AddDevice struct {
	SerialNumber string             `json:"serial_number"`
	Model        string             `json:"model"`
	PetID        primitive.ObjectID `json:"pet_id"`
	Status       string             `json:"status"`
	LastSyncTime int64              `json:"last_sync_time"`
	CreatedAt    time.Time          `json:"created_at"`
}

func NewDevice(r *http.Request) (*AddDevice, error) {
	bodyReader := r.Body
	if bodyReader == nil {
		return nil, errors.New("missing body")
	}

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	var device AddDevice
	err = json.Unmarshal(body, &device)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(device.SerialNumber) == "" {
		return nil, errors.New("serial_number is required")
	}

	if strings.TrimSpace(device.Model) == "" {
		return nil, errors.New("model is required")
	}

	return &device, nil
}
