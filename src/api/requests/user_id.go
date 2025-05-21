package requests

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type UserID struct {
	ID primitive.ObjectID `json:"_id"`
}

func NewUserID(r *http.Request) (*UserID, error) {
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		return nil, errors.New("missing pet ID in URL")
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, err
	}

	return &UserID{
		ID: id,
	}, nil
}
