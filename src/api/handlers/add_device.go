package handlers

import (
	"encoding/json"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"
)

func AddDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req, err := requests.NewDevice(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format: " + err.Error()})
		return
	}

	if strings.TrimSpace(req.SerialNumber) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Serial number is required"})
		return
	}

	if strings.TrimSpace(req.Model) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Model is required"})
		return
	}

	if strings.TrimSpace(req.Status) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Device status is required"})
		return
	}

	validStatuses := map[string]bool{
		"active":   true,
		"inactive": true,
		"offline":  true,
	}

	if !validStatuses[req.Status] {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid device status. Must be active, inactive, or offline"})
		return
	}

	devicesDB := MongoDB(r).Devices()
	existingDevices, err := devicesDB.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to check existing devices"})
		return
	}

	for _, d := range existingDevices {
		if d.SerialNumber == req.SerialNumber {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Device with this serial number already exists"})
			return
		}
	}

	var lastSyncTime primitive.Timestamp
	if req.LastSyncTime > 0 {
		lastSyncTime = primitive.Timestamp{
			T: uint32(req.LastSyncTime),
			I: 0,
		}
	} else {
		currentTime := time.Now().Unix()
		lastSyncTime = primitive.Timestamp{
			T: uint32(currentTime),
			I: 0,
		}
	}

	device := data.Device{
		ID:           primitive.NewObjectID(),
		SerialNumber: req.SerialNumber,
		Model:        req.Model,
		Status:       req.Status,
		LastSyncTime: lastSyncTime,
		CreatedAt:    time.Now(),
	}

	if !req.PetID.IsZero() {
		petDB := MongoDB(r).Pets()
		pet, err := petDB.Get(req.PetID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Pet not found"})
			return
		}
		device.PetID = pet.ID
	}

	err = devicesDB.Insert(&device)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to add device: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Device added successfully",
		"deviceID": device.ID.Hex(),
	})
}
