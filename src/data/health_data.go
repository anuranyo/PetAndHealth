package data

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HealthDataDB interface {
	Get(pet_id primitive.ObjectID) (*HealthData, error)
	Insert(*HealthData) error
	Update(pet_id primitive.ObjectID, updateFields bson.M) error
	GetByFilter(filter bson.M) ([]*HealthData, error)
}

type HealthData struct {
	ID          primitive.ObjectID  `bson:"_id"`
	PetID       primitive.ObjectID  `bson:"pet_id"`
	Activity    float64             `bson:"activity"`
	SleepHours  float64             `bson:"sleep_hours"`
	Temperature float64             `bson:"temperature"`
	Time        primitive.Timestamp `bson:"time"`
}
