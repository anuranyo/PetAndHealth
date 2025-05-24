package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Константы для нормальных значений здоровья животных
const (
	// Температура (°C)
	NormalTemperatureMin = 37.5
	NormalTemperatureMax = 39.5
	NormalTemperatureAvg = 38.5

	// Сон (часы)
	NormalSleepMin = 8.0
	NormalSleepMax = 16.0
	NormalSleepAvg = 12.0
)

// HealthStatus представляет статус здоровья питомца
type HealthStatus struct {
	PetID             string    `json:"pet_id"`
	PetName           string    `json:"pet_name"`
	PetSpecies        string    `json:"pet_species"`
	LastCheckTime     time.Time `json:"last_check_time"`
	TemperatureStatus string    `json:"temperature_status"`
	TemperatureValue  float64   `json:"temperature_value"`
	SleepStatus       string    `json:"sleep_status"`
	SleepValue        float64   `json:"sleep_value"`
	OverallStatus     string    `json:"overall_status"`
	Issues            []string  `json:"issues,omitempty"`
}

// OwnerHealthSummary представляет общую сводку здоровья питомцев владельца
type OwnerHealthSummary struct {
	OwnerID       string         `json:"owner_id"`
	TotalPets     int            `json:"total_pets"`
	HealthyPets   int            `json:"healthy_pets"`
	ProblemsCount int            `json:"problems_count"`
	LastUpdated   time.Time      `json:"last_updated"`
	PetsHealth    []HealthStatus `json:"pets_health"`
}

// GetHealthPetsByOwner проверяет здоровье всех питомцев владельца
func GetHealthPetsByOwner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получаем ID текущего пользователя
	currentUserID, ok := r.Context().Value(UserIDContextKey).(primitive.ObjectID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	currentUserRole, _ := r.Context().Value(UserRoleContextKey).(string)

	// Для владельцев проверяем только их питомцев
	var ownerID primitive.ObjectID
	if currentUserRole == "user" {
		ownerID = currentUserID
	} else {
		// Для админов и ветов можно передать owner_id как query параметр
		ownerIDStr := r.URL.Query().Get("owner_id")
		if ownerIDStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "owner_id parameter is required"})
			return
		}

		var err error
		ownerID, err = primitive.ObjectIDFromHex(ownerIDStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid owner_id format"})
			return
		}
	}

	// Получаем всех питомцев владельца
	petsDB := MongoDB(r).Pets()
	filter := bson.M{"owner_id": ownerID}

	pets, err := petsDB.GetByFilter(filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to retrieve pets"})
		return
	}

	if len(pets) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OwnerHealthSummary{
			OwnerID:       ownerID.Hex(),
			TotalPets:     0,
			HealthyPets:   0,
			ProblemsCount: 0,
			LastUpdated:   time.Now(),
			PetsHealth:    []HealthStatus{},
		})
		return
	}

	// Проверяем здоровье каждого питомца
	healthDataDB := MongoDB(r).HealthData()
	var petsHealth []HealthStatus
	healthyCount := 0
	problemsCount := 0

	for _, pet := range pets {
		// Получаем последние данные здоровья питомца
		filter := bson.M{"pet_id": pet.ID}

		allHealthData, err := healthDataDB.GetByFilter(filter)
		if err != nil || len(allHealthData) == 0 {
			// Если нет данных о здоровье
			petsHealth = append(petsHealth, HealthStatus{
				PetID:         pet.ID.Hex(),
				PetName:       pet.Name,
				PetSpecies:    pet.Species,
				LastCheckTime: time.Time{},
				OverallStatus: "no_data",
				Issues:        []string{"No health data available"},
			})
			problemsCount++
			continue
		}

		// Находим самую последнюю запись
		var latestData *data.HealthData
		for _, healthRecord := range allHealthData {
			if latestData == nil || healthRecord.Time.T > latestData.Time.T {
				latestData = healthRecord
			}
		}

		if latestData == nil {
			petsHealth = append(petsHealth, HealthStatus{
				PetID:         pet.ID.Hex(),
				PetName:       pet.Name,
				PetSpecies:    pet.Species,
				LastCheckTime: time.Time{},
				OverallStatus: "no_data",
				Issues:        []string{"No health data available"},
			})
			problemsCount++
			continue
		}

		healthStatus := analyzeHealthData(pet, latestData)
		petsHealth = append(petsHealth, healthStatus)

		if healthStatus.OverallStatus == "healthy" {
			healthyCount++
		} else {
			problemsCount++
		}
	}

	// Формируем ответ
	summary := OwnerHealthSummary{
		OwnerID:       ownerID.Hex(),
		TotalPets:     len(pets),
		HealthyPets:   healthyCount,
		ProblemsCount: problemsCount,
		LastUpdated:   time.Now(),
		PetsHealth:    petsHealth,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summary)
}

// analyzeHealthData анализирует данные здоровья и возвращает статус
func analyzeHealthData(pet *data.Pet, healthData *data.HealthData) HealthStatus {
	status := HealthStatus{
		PetID:            pet.ID.Hex(),
		PetName:          pet.Name,
		PetSpecies:       pet.Species,
		LastCheckTime:    time.Unix(int64(healthData.Time.T), 0),
		TemperatureValue: healthData.Temperature,
		SleepValue:       healthData.SleepHours,
		Issues:           []string{},
	}

	// Анализ температуры
	if healthData.Temperature >= NormalTemperatureMin && healthData.Temperature <= NormalTemperatureMax {
		status.TemperatureStatus = "normal"
	} else if healthData.Temperature < NormalTemperatureMin {
		status.TemperatureStatus = "low"
		status.Issues = append(status.Issues, "Temperature too low")
	} else {
		status.TemperatureStatus = "high"
		status.Issues = append(status.Issues, "Temperature too high")
	}

	// Анализ сна
	if healthData.SleepHours >= NormalSleepMin && healthData.SleepHours <= NormalSleepMax {
		status.SleepStatus = "normal"
	} else if healthData.SleepHours < NormalSleepMin {
		status.SleepStatus = "insufficient"
		status.Issues = append(status.Issues, "Insufficient sleep")
	} else {
		status.SleepStatus = "excessive"
		status.Issues = append(status.Issues, "Excessive sleep")
	}

	// Определяем общий статус
	if len(status.Issues) == 0 {
		status.OverallStatus = "healthy"
	} else if len(status.Issues) == 1 {
		status.OverallStatus = "minor_issues"
	} else {
		status.OverallStatus = "attention_needed"
	}

	return status
}
