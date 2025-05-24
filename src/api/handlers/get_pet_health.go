package handlers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetPetHealth возвращает все данные о здоровье конкретного питомца
func GetPetHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получаем ID питомца из URL параметров
	req, err := requests.NewPetID(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid pet ID format"})
		return
	}

	if req.ID.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Pet ID is required"})
		return
	}

	// Проверяем существование питомца
	petsDB := MongoDB(r).Pets()
	pet, err := petsDB.Get(req.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Pet not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to retrieve pet information"})
		return
	}

	// Проверяем права доступа
	currentUserID, ok := r.Context().Value(UserIDContextKey).(primitive.ObjectID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	currentUserRole, _ := r.Context().Value(UserRoleContextKey).(string)

	// Владелец может видеть только своих питомцев, админ и vet - всех
	if currentUserRole == "user" && pet.OwnerID != currentUserID {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "You can only view your own pets' health data"})
		return
	}

	// Получаем все данные о здоровье питомца
	healthDataDB := MongoDB(r).HealthData()
	filter := bson.M{"pet_id": req.ID}

	healthData, err := healthDataDB.GetByFilter(filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to retrieve health data"})
		return
	}

	// Сортируем по времени (новые записи сначала)
	sort.Slice(healthData, func(i, j int) bool {
		return healthData[i].Time.T > healthData[j].Time.T
	})

	// Формируем ответ с информацией о питомце и его данными здоровья
	response := struct {
		Pet        interface{} `json:"pet"`
		HealthData interface{} `json:"health_data"`
		Count      int         `json:"count"`
	}{
		Pet: struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Species string `json:"species"`
			Breed   string `json:"breed"`
			Age     int    `json:"age"`
		}{
			ID:      pet.ID.Hex(),
			Name:    pet.Name,
			Species: pet.Species,
			Breed:   pet.Breed,
			Age:     pet.Age,
		},
		HealthData: healthData,
		Count:      len(healthData),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
