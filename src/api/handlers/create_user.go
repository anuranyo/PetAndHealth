package handlers

import (
	"encoding/json"
	"errors"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type CreateUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func NewCreateUser(r *http.Request) (*CreateUserRequest, error) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req, err := NewCreateUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format"})
		return
	}

	if err := validateUsersData(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	validRoles := map[string]bool{
		"admin": true,
		"user":  true,
		"vet":   true,
	}
	if !validRoles[req.Role] {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid role. Must be 'admin', 'user', or 'vet'"})
		return
	}

	usersDB := MongoDB(r).Users()
	existingUser, _ := usersDB.FindByEmail(req.Email)
	if existingUser != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "User with this email already exists"})
		return
	}

	password, err := utils.GenerateRandomPassword(12)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to generate password"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to process password"})
		return
	}

	user := data.User{
		ID:           primitive.NewObjectID(),
		FullName:     req.FullName,
		Email:        req.Email,
		Role:         req.Role,
		PasswordHash: string(hashedPassword),
		PetsID:       []primitive.ObjectID{},
	}

	err = usersDB.Insert(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create user"})
		return
	}

	roleName := map[string]string{
		"admin": "Administrator",
		"user":  "User",
		"vet":   "Veterinarian",
	}[req.Role]

	err = utils.SendPasswordEmail(user.Email, user.FullName, password, roleName)
	emailStatus := "Email with login details sent successfully"

	if err != nil {
		log.Printf("Failed to send password email: %v", err)
		emailStatus = "User created, but failed to send email with login details"
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":     roleName + " account created successfully",
		"adminID":     user.ID.Hex(),
		"role":        user.Role,
		"emailStatus": emailStatus,
	})
}

func validateUsersData(req *CreateUserRequest) error {
	if strings.TrimSpace(req.FullName) == "" {
		return errors.New("Full name is required")
	}

	if strings.TrimSpace(req.Email) == "" {
		return errors.New("Email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("Email is invalid")
	}

	return nil
}
