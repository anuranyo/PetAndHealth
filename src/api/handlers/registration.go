package handlers

import (
	"encoding/json"
	"errors"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/middle"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
)

type RegistrationResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
	Token   string `json:"token"`
}

func Registration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req, err := requests.NewRegistration(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format"})
		return
	}

	if err := validateUserData(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	usersDB := MongoDB(r).Users()
	existingUser, _ := usersDB.FindByEmail(req.Email)
	if existingUser != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "User with this email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to process password"})
		return
	}

	user := data.User{
		ID:           primitive.NewObjectID(),
		FullName:     req.FullName,
		Email:        strings.ToLower(req.Email),
		Role:         "user",
		PasswordHash: string(hashedPassword),
		PetsID:       []primitive.ObjectID{},
	}

	err = usersDB.Insert(&user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "User with this email already exists"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to save user"})
		return
	}

	token, err := middle.GenerateJWT(user.ID, user.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User registered successfully, but token generation failed",
			"userID":  user.ID.Hex(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := RegistrationResponse{
		Message: "User registered successfully",
		UserID:  user.ID.Hex(),
		Token:   token,
	}
	json.NewEncoder(w).Encode(response)
}

func validateUserData(req *requests.Registration) error {
	if strings.TrimSpace(req.FullName) == "" {
		return errors.New("Full name is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("Email is invalid")
	}

	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}
