package data

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UsersDB interface {
	Get(id primitive.ObjectID) (*User, error)
	Insert(*User) error
	Update(id primitive.ObjectID, updateFields bson.M) error
	FindByEmail(email string) (*User, error)
	GetAll() ([]*User, error)
	UpdatePets(userID primitive.ObjectID, petID primitive.ObjectID) error
	Delete(id primitive.ObjectID) error
	ResetPassword(email string, newPasswordHash string) error
	RequirePasswordChange(id primitive.ObjectID, required bool) error
}

type User struct {
	ID                     primitive.ObjectID   `bson:"_id"`
	FullName               string               `bson:"full_name"`
	Role                   string               `bson:"role"`
	Email                  string               `bson:"email"`
	PasswordHash           string               `bson:"password_hash"`
	PetsID                 []primitive.ObjectID `bson:"pets_id"`
	PasswordResetTime      *time.Time           `bson:"password_reset_time,omitempty"`
	PasswordChangeRequired bool                 `bson:"password_change_required,omitempty"`
}
