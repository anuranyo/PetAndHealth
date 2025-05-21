package handlers

import (
	"encoding/json"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func UpdatePet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	petIDStr := chi.URLParam(r, "id")
	if petIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Pet ID is required"})
		return
	}

	petID, err := primitive.ObjectIDFromHex(petIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid pet ID format"})
		return
	}

	req, err := requests.NewUpdatePet(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format"})
		return
	}

	petsDB := MongoDB(r).Pets()
	usersDB := MongoDB(r).Users()

	currentPet, err := petsDB.Get(petID)
	if err != nil || currentPet == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Pet not found"})
		return
	}

	updateFields := bson.M{}

	if req.Name != "" {
		updateFields["name"] = req.Name
	}
	if req.Species != "" {
		updateFields["species"] = req.Species
	}
	if req.Breed != "" {
		updateFields["breed"] = req.Breed
	}
	if req.Age > 0 {
		updateFields["age"] = req.Age
	}

	if !req.OwnerID.IsZero() && req.OwnerID != currentPet.OwnerID {
		if !currentPet.OwnerID.IsZero() {
			oldOwner, err := usersDB.Get(currentPet.OwnerID)
			if err == nil && oldOwner != nil {
				newPetsArray := []primitive.ObjectID{}
				for _, pID := range oldOwner.PetsID {
					if pID != petID {
						newPetsArray = append(newPetsArray, pID)
					}
				}

				err = usersDB.Update(currentPet.OwnerID, bson.M{
					"$set": bson.M{"pets_id": newPetsArray},
				})
				if err != nil {
					log.Printf("Error removing pet from old owner: %v", err)
				}
			}
		}

		newOwner, err := usersDB.Get(req.OwnerID)
		if err == nil && newOwner != nil {
			alreadyExists := false
			for _, pID := range newOwner.PetsID {
				if pID == petID {
					alreadyExists = true
					break
				}
			}

			if !alreadyExists {
				err = usersDB.UpdatePets(req.OwnerID, petID)
				if err != nil {
					log.Printf("Error adding pet to new owner: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update owner's pet list"})
					return
				}
			}
		}

		updateFields["owner_id"] = req.OwnerID
	}

	if len(updateFields) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No fields to update"})
		return
	}

	err = petsDB.Update(petID, updateFields)
	if err != nil {
		log.Printf("Error updating pet: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update pet: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Pet updated successfully",
		"petID":   petID.Hex(),
	})
}
