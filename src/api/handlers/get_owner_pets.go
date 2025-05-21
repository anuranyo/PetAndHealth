package handlers

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOwnerPets(w http.ResponseWriter, r *http.Request) {
	ownerIDHex := r.URL.Query().Get("owner_id")
	if ownerIDHex == "" {
		http.Error(w, "Missing owner_id parameter", http.StatusBadRequest)
		return
	}

	ownerID, err := primitive.ObjectIDFromHex(ownerIDHex)
	if err != nil {
		http.Error(w, "Invalid owner_id format", http.StatusBadRequest)
		return
	}

	usersDB := MongoDB(r).Users()
	petsDB := MongoDB(r).Pets()

	user, err := usersDB.Get(ownerID)
	if err != nil {
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	if len(user.PetsID) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	// Знаходимо тварин за їхніми ID
	filter := bson.M{"_id": bson.M{"$in": user.PetsID}}
	pets, err := petsDB.GetByFilter(filter)
	if err != nil {
		http.Error(w, "Failed to retrieve pets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pets)
}
