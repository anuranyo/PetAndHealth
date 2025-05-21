package handlers

import (
	"encoding/json"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr := chi.URLParam(r, "id")
	if userIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "User ID is required"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid user ID format"})
		return
	}

	req, err := requests.NewUpdateUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format"})
		return
	}

	req.ID = userID

	currentUserID, ok := r.Context().Value(UserIDContextKey).(primitive.ObjectID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	currentUserRole, _ := r.Context().Value(UserRoleContextKey).(string)

	if currentUserRole != "admin" && userID != currentUserID {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "You can only update your own profile"})
		return
	}

	usersDB := MongoDB(r).Users()
	currentUser, err := usersDB.Get(userID)
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
		if existingUser != nil && existingUser.ID != userID {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Email already in use"})
			return
		}
		updateFields["email"] = strings.ToLower(req.Email)
	}

	if req.Role != "" && currentUserRole == "admin" {
		validRoles := map[string]bool{
			"admin": true,
			"user":  true,
			"vet":   true,
		}
		if !validRoles[req.Role] {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid role"})
			return
		}
		updateFields["role"] = req.Role
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

	err = usersDB.Update(userID, bson.M{"$set": updateFields})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User updated successfully",
		"userID":  userID.Hex(),
	})
}
