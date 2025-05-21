package requests

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type PetReportTimeRange struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

func NewPetReportTimeRange(r *http.Request) (*PetReportTimeRange, error) {
	bodyReader := r.Body
	if bodyReader == nil {
		return nil, errors.New("missing body")
	}

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	var timeRange PetReportTimeRange
	err = json.Unmarshal(body, &timeRange)
	if err != nil {
		return nil, err
	}

	return &timeRange, nil
}
