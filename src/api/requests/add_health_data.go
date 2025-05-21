package requests

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
)

type AddHealthData struct {
	PetID       primitive.ObjectID `json:"pet_id"`
	Activity    float64            `json:"activity"`
	SleepHours  float64            `json:"sleep_hours"`
	Temperature float64            `json:"temperature"`
}

func NewHealthData(r *http.Request) (*AddHealthData, error) {
	bodyReader := r.Body
	if bodyReader == nil {
		return nil, errors.New("missing body")
	}

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	var data AddHealthData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
