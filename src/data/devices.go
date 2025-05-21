package data

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DevicesDB interface {
	Get(id primitive.ObjectID) (*Device, error)
	Insert(*Device) error
	Update(id primitive.ObjectID, updateFields bson.M) error
	GetAll() ([]*Device, error)
}

type Device struct {
	ID           primitive.ObjectID  `bson:"_id"`
	SerialNumber string              `bson:"serial_number"`
	Model        string              `bson:"model"`
	PetID        primitive.ObjectID  `bson:"pet_id"`
	Status       string              `bson:"status"`
	LastSyncTime primitive.Timestamp `bson:"last_sync_time"`
	CreatedAt    time.Time           `bson:"created_at"`
}
