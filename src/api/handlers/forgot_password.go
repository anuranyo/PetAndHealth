package handlers

import (
	"encoding/json"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format"})
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Email is required"})
		return
	}

	usersDB := MongoDB(r).Users()
	user, err := usersDB.FindByEmail(strings.ToLower(req.Email))

	if err != nil {
		log.Printf("Password reset requested for non-existent email: %s, error: %v", req.Email, err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "If this email exists in our system, password reset instructions will be sent",
		})
		return
	}

	newPassword, err := utils.GenerateRandomPassword(12)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to process password reset"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to process password"})
		return
	}

	err = usersDB.Update(user.ID, bson.M{
		"$set": bson.M{
			"password_hash":     string(hashedPassword),
			"password_reset_at": time.Now(),
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update password"})
		return
	}

	roleName := map[string]string{
		"admin": "Administrator",
		"user":  "User",
		"vet":   "Veterinarian",
	}[user.Role]

	err = utils.SendPasswordResetEmail(user.Email, user.FullName, newPassword, roleName)
	if err != nil {
		log.Printf("Failed to send password reset email to %s: %v", user.Email, err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If this email exists in our system, password reset instructions will be sent",
	})
}
