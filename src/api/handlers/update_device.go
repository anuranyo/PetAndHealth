package handlers

import (
	"encoding/json"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func UpdateDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Отримуємо ID пристрою з URL параметрів
	deviceIDStr := chi.URLParam(r, "id")
	if deviceIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Device ID is required"})
		return
	}

	// Конвертуємо ID в ObjectID
	deviceID, err := primitive.ObjectIDFromHex(deviceIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid device ID format"})
		return
	}

	// Парсимо тіло запиту
	req, err := requests.NewUpdateDevice(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format: " + err.Error()})
		return
	}

	// Встановлюємо ID з URL параметра
	req.ID = deviceID

	devicesDB := MongoDB(r).Devices()
	existingDevice, err := devicesDB.Get(deviceID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Device not found"})
		return
	}

	if req.SerialNumber != "" && req.SerialNumber != existingDevice.SerialNumber {
		devices, err := devicesDB.GetAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to check existing devices"})
			return
		}

		for _, d := range devices {
			if d.ID != deviceID && d.SerialNumber == req.SerialNumber {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Device with this serial number already exists"})
				return
			}
		}
	}

	updateFields := bson.M{}

	if req.SerialNumber != "" {
		updateFields["serial_number"] = req.SerialNumber
	}

	if req.Model != "" {
		updateFields["model"] = req.Model
	}

	if req.Status != "" {
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

		updateFields["status"] = req.Status
	}

	if !req.PetID.IsZero() {
		petDB := MongoDB(r).Pets()
		pet, err := petDB.Get(req.PetID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Pet not found"})
			return
		}

		updateFields["pet_id"] = pet.ID
	}

	if req.LastSyncTime > 0 {
		lastSyncTime := primitive.Timestamp{
			T: uint32(req.LastSyncTime),
			I: 0,
		}
		updateFields["last_sync_time"] = lastSyncTime
	}

	if len(updateFields) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No fields to update"})
		return
	}

	err = devicesDB.Update(deviceID, updateFields)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update device: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Device updated successfully",
		"deviceID": deviceID.Hex(),
	})
}
