package handlers

import (
	"encoding/json"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"net/http"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.ID.IsZero() {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	usersDB := MongoDB(r).Users()

	err = usersDB.Delete(req.ID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted successfully",
	})
}
