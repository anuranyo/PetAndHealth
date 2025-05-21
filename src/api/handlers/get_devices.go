package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func GetDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	devicesDB := MongoDB(r).Devices()
	devices, err := devicesDB.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to fetch devices"})
		return
	}

	petsDB := MongoDB(r).Pets()
	usersDB := MongoDB(r).Users()

	type DeviceResponse struct {
		ID           string    `json:"ID"`
		SerialNumber string    `json:"SerialNumber"`
		Model        string    `json:"Model"`
		Status       string    `json:"Status"`
		PetID        string    `json:"PetID,omitempty"`
		PetName      string    `json:"PetName,omitempty"`
		PetOwner     string    `json:"PetOwner,omitempty"`
		LastSyncTime int64     `json:"LastSyncTime"`
		CreatedAt    time.Time `json:"CreatedAt"`
	}

	var response []DeviceResponse

	for _, device := range devices {
		deviceResp := DeviceResponse{
			ID:           device.ID.Hex(),
			SerialNumber: device.SerialNumber,
			Model:        device.Model,
			Status:       device.Status,
			LastSyncTime: int64(device.LastSyncTime.T),
			CreatedAt:    device.CreatedAt,
		}

		if !device.PetID.IsZero() {
			deviceResp.PetID = device.PetID.Hex()

			pet, err := petsDB.Get(device.PetID)
			if err == nil {
				deviceResp.PetName = pet.Name

				if !pet.OwnerID.IsZero() {
					owner, err := usersDB.Get(pet.OwnerID)
					if err == nil {
						deviceResp.PetOwner = owner.FullName
					}
				}
			}
		}

		response = append(response, deviceResp)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
