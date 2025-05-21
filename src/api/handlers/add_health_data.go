package handlers

import (
	"encoding/json"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"
)

func AddHealthData(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewHealthData(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currentTime := time.Now()
	lastSyncTime := primitive.Timestamp{
		T: uint32(currentTime.Unix()),
		I: 0,
	}

	healthData := data.HealthData{
		ID:          primitive.NewObjectID(),
		PetID:       req.PetID,
		Activity:    req.Activity,
		SleepHours:  req.SleepHours,
		Temperature: req.Temperature,
		Time:        lastSyncTime,
	}

	healthDataDB := MongoDB(r).HealthData()

	err = healthDataDB.Insert(&healthData)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			w.WriteHeader(http.StatusConflict)
			return
		}
		http.Error(w, "Failed to add data", http.StatusInternalServerError)
		return
	}

	petDB := MongoDB(r).Pets()
	pet, err := petDB.Get(req.PetID)
	if err != nil || pet == nil {
		http.Error(w, "Failed to fetch pet information", http.StatusInternalServerError)
		return
	}

	ownerID := pet.OwnerID

	notificationsDB := MongoDB(r).Notifications()
	utils.CheckPetHealthAndNotify(notificationsDB, ownerID, req.Temperature, req.SleepHours, currentTime)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Data added successfully",
		"healthDataID": healthData.ID.Hex(),
		"time":         currentTime.Format(time.RFC3339),
	})
}
