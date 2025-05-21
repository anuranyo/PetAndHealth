package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func UserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userID primitive.ObjectID

	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		currentUserID, ok := r.Context().Value(UserIDContextKey).(primitive.ObjectID)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
			return
		}
		userID = currentUserID
	} else {
		var err error
		userID, err = primitive.ObjectIDFromHex(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid user ID format"})
			return
		}

		currentUserRole, _ := r.Context().Value(UserRoleContextKey).(string)
		currentUserID, _ := r.Context().Value(UserIDContextKey).(primitive.ObjectID)

		if currentUserRole != "admin" && currentUserID != userID {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "You can only view your own profile"})
			return
		}
	}

	usersDB := MongoDB(r).Users()
	user, err := usersDB.Get(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "User not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to retrieve user information"})
		return
	}

	user.PasswordHash = ""

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
