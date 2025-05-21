package handlers

import (
	"encoding/json"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req, err := requests.NewUpdateUserProfile(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format"})
		return
	}

	currentUserID, ok := r.Context().Value(UserIDContextKey).(primitive.ObjectID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	usersDB := MongoDB(r).Users()
	currentUser, err := usersDB.Get(currentUserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "User not found"})
		return
	}

	updateFields := bson.M{}

	if strings.TrimSpace(req.FullName) != "" {
		updateFields["full_name"] = req.FullName
	}

	if strings.TrimSpace(req.Email) != "" && req.Email != currentUser.Email {
		existingUser, _ := usersDB.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != currentUserID {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Email already in use"})
			return
		}
		updateFields["email"] = strings.ToLower(req.Email)
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to process password"})
			return
		}
		updateFields["password_hash"] = string(hashedPassword)
	}

	if len(updateFields) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No fields to update"})
		return
	}

	err = usersDB.Update(currentUserID, bson.M{"$set": updateFields})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update user profile"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile updated successfully",
		"userID":  currentUserID.Hex(),
	})
}
