package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func AssignDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	deviceIDStr := chi.URLParam(r, "id")
	if deviceIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Device ID is required"})
		return
	}

	deviceID, err := primitive.ObjectIDFromHex(deviceIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid device ID format"})
		return
	}

	var req struct {
		PetID string `json:"pet_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format"})
		return
	}

	petID, err := primitive.ObjectIDFromHex(req.PetID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid pet ID format"})
		return
	}

	devicesDB := MongoDB(r).Devices()
	device, err := devicesDB.Get(deviceID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Device not found"})
		return
	}

	if !device.PetID.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Device is already assigned to pet"})
		return
	}

	petsDB := MongoDB(r).Pets()
	pet, err := petsDB.Get(petID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Pet not found"})
		return
	}

	updateFields := bson.M{
		"pet_id": pet.ID,
		"status": "active",
	}

	err = devicesDB.Update(deviceID, updateFields)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to assign device"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Device assigned successfully",
		"deviceID": deviceID.Hex(),
		"petID":    petID.Hex(),
	})
}

func UnassignDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	deviceIDStr := chi.URLParam(r, "id")
	if deviceIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Device ID is required"})
		return
	}

	deviceID, err := primitive.ObjectIDFromHex(deviceIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid device ID format"})
		return
	}

	devicesDB := MongoDB(r).Devices()
	device, err := devicesDB.Get(deviceID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Device not found"})
		return
	}

	if device.PetID.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Device is not assigned to any pet"})
		return
	}

	updateFields := bson.M{
		"pet_id": primitive.NilObjectID,
		"status": "inactive",
	}

	err = devicesDB.Update(deviceID, updateFields)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to unassign device"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Device unassigned successfully",
		"deviceID": deviceID.Hex(),
	})
}
