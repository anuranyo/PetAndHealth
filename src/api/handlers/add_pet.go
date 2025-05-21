package handlers

import (
	"encoding/json"
	"errors"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"
)

func AddPet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req, err := requests.NewPet(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format: " + err.Error()})
		return
	}

	if err := validatePetData(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	usersDB := MongoDB(r).Users()
	owner, err := usersDB.Get(req.OwnerID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Pet owner not found"})
		return
	}

	pet := data.Pet{
		ID:      primitive.NewObjectID(),
		Name:    req.Name,
		Species: req.Species,
		Breed:   req.Breed,
		Age:     req.Age,
		OwnerID: owner.ID,
	}

	petsDB := MongoDB(r).Pets()
	err = petsDB.Insert(&pet)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Pet already exists"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to add pet"})
		return
	}

	err = usersDB.UpdatePets(req.OwnerID, pet.ID)
	if err != nil {
		_ = petsDB.Delete(pet.ID)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update user with pet ID"})
		return
	}

	notification := data.Notification{
		ID:      primitive.NewObjectID(),
		UserID:  req.OwnerID,
		Title:   "New pet added",
		Message: "Your pet " + pet.Name + " has been added to your profile.",
		Type:    "info",
		Time:    time.Now(),
		Read:    false,
	}

	notificationsDB := MongoDB(r).Notifications()
	_ = notificationsDB.Insert(&notification)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Pet added successfully",
		"petID":   pet.ID.Hex(),
	})
}

func validatePetData(pet *requests.AddPet) error {
	if strings.TrimSpace(pet.Name) == "" {
		return errors.New("pet name is required")
	}

	if strings.TrimSpace(pet.Species) == "" {
		return errors.New("pet species is required")
	}

	if pet.Age < 0 || pet.Age > 100 {
		return errors.New("invalid pet age")
	}

	if pet.OwnerID.IsZero() {
		return errors.New("owner ID is required")
	}

	return nil
}
