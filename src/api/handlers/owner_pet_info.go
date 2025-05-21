package handlers

import (
	"encoding/json"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func OwnerPetInfo(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewPetID(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.ID.IsZero() {
		http.Error(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	petsDB := MongoDB(r).Pets()
	petWithHealth, err := petsDB.Get(req.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Pet not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve pet information", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(petWithHealth)
}
