package handlers

import (
	"encoding/json"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/middle"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type AuthResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
	Role    string `json:"role"`
	Token   string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req, err := requests.NewAuth(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format"})
		return
	}

	usersDB := MongoDB(r).Users()
	user, err := usersDB.FindByEmail(req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "User not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to authenticate user"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid password"})
		return
	}

	token, err := middle.GenerateJWT(user.ID, user.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to generate authentication token"})
		return
	}

	w.WriteHeader(http.StatusOK)
	response := AuthResponse{
		Message: "User authenticated successfully",
		UserID:  user.ID.Hex(),
		Role:    user.Role,
		Token:   token,
	}
	json.NewEncoder(w).Encode(response)
}
