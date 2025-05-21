package data

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type NotificationsDB interface {
	Get(id string) (*Notification, error)
	Insert(*Notification) error
	Update(id string, updateFields bson.M) error
	GetForUser(userID primitive.ObjectID) ([]*Notification, error)
	MarkAsRead(id primitive.ObjectID) error
	GetUnreadCount(userID primitive.ObjectID) (int, error)
}

type Notification struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Title     string             `bson:"title"`
	Message   string             `bson:"message"`
	Type      string             `bson:"type"`
	Time      time.Time          `bson:"time"`
	Read      bool               `bson:"read"`
	Delivered bool               `bson:"delivered"`
}
