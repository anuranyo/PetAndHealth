package handlers

import (
	"encoding/json"
	"net/http"
)

func GetPets(w http.ResponseWriter, r *http.Request) {
	petsDB := MongoDB(r).Pets()

	pets, err := petsDB.GetAll()
	if err != nil {
		http.Error(w, "Failed to retrieve pets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(pets)
}
