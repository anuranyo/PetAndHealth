package handlers

import (
	"encoding/json"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	usersDB := MongoDB(r).Users()

	users, err := usersDB.GetAll()
	if err != nil {
		http.Error(w, "Failed to retrieve pets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(users)
}
