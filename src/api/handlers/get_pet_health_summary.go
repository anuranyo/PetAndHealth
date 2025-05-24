package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PetHealthSummary представляет краткую сводку здоровья питомца для уведомлений
type PetHealthSummary struct {
	PetID         string    `json:"pet_id"`
	PetName       string    `json:"pet_name"`
	PetSpecies    string    `json:"pet_species"`
	PetBreed      string    `json:"pet_breed"`
	PetAge        int       `json:"pet_age"`
	LastCheckTime time.Time `json:"last_check_time"`
	OverallStatus string    `json:"overall_status"` // "healthy", "minor_issues", "attention_needed", "no_data"
	HealthScore   int       `json:"health_score"`   // 0-100
	AlertLevel    string    `json:"alert_level"`    // "none", "low", "medium", "high"

	// Детальные показатели
	Temperature HealthMetric `json:"temperature"`
	Sleep       HealthMetric `json:"sleep"`
	Activity    HealthMetric `json:"activity"`

	// Рекомендации и предупреждения
	Issues          []string `json:"issues,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`

	// Для уведомлений
	RequiresAttention bool   `json:"requires_attention"`
	NotificationLevel string `json:"notification_level"` // "info", "warning", "urgent"
	DataAge           string `json:"data_age"`           // "recent", "outdated", "none"
}

// HealthMetric представляет метрику здоровья с оценкой
type HealthMetric struct {
	Value  float64 `json:"value"`
	Status string  `json:"status"` // "normal", "low", "high", "critical"
	Score  int     `json:"score"`  // 0-100
	Normal string  `json:"normal"` // нормальный диапазон для отображения
}

// GetPetHealthSummary возвращает краткую сводку здоровья конкретного питомца владельца
func GetPetHealthSummary(w http.ResponseWriter, r *http.Request) {
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

	// Получаем ID текущего пользователя
	currentUserID, ok := r.Context().Value(UserIDContextKey).(primitive.ObjectID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Проверяем существование питомца и права доступа
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

	// Проверяем, что питомец принадлежит текущему пользователю
	if pet.OwnerID != currentUserID {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "You can only view your own pets' health summary"})
		return
	}

	// Получаем данные здоровья питомца
	healthDataDB := MongoDB(r).HealthData()
	filter := bson.M{"pet_id": req.ID}

	allHealthData, err := healthDataDB.GetByFilter(filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to retrieve health data"})
		return
	}

	// Формируем сводку здоровья
	summary := createPetHealthSummary(pet, allHealthData)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summary)
}

// createPetHealthSummary создает краткую сводку здоровья питомца
func createPetHealthSummary(pet *data.Pet, healthData []*data.HealthData) PetHealthSummary {
	summary := PetHealthSummary{
		PetID:             pet.ID.Hex(),
		PetName:           pet.Name,
		PetSpecies:        pet.Species,
		PetBreed:          pet.Breed,
		PetAge:            pet.Age,
		Issues:            []string{},
		Recommendations:   []string{},
		RequiresAttention: false,
	}

	// Если нет данных о здоровье
	if len(healthData) == 0 {
		summary.OverallStatus = "no_data"
		summary.HealthScore = 0
		summary.AlertLevel = "medium"
		summary.NotificationLevel = "warning"
		summary.DataAge = "none"
		summary.Issues = append(summary.Issues, "No health data available")
		summary.Recommendations = append(summary.Recommendations, "Check device connection and ensure regular monitoring")
		summary.RequiresAttention = true
		return summary
	}

	// Находим самую последнюю запись
	var latestData *data.HealthData
	for _, healthRecord := range healthData {
		healthRecordTime := time.Unix(int64(healthRecord.Time.T), 0)
		var latestDataTime time.Time
		if latestData != nil {
			latestDataTime = time.Unix(int64(latestData.Time.T), 0)
		}
		if latestData == nil || healthRecordTime.After(latestDataTime) {
			latestData = healthRecord
		}
	}

	if latestData == nil {
		summary.OverallStatus = "no_data"
		summary.HealthScore = 0
		summary.AlertLevel = "medium"
		summary.NotificationLevel = "warning"
		summary.DataAge = "none"
		summary.RequiresAttention = true
		return summary
	}

	summary.LastCheckTime = time.Unix(int64(latestData.Time.T), 0)

	// Определяем возраст данных
	dataAge := time.Since(time.Unix(int64(latestData.Time.T), 0))
	if dataAge <= 1*time.Hour {
		summary.DataAge = "recent"
	} else if dataAge <= 24*time.Hour {
		summary.DataAge = "outdated"
	} else {
		summary.DataAge = "old"
		summary.Issues = append(summary.Issues, "Health data is outdated")
		summary.Recommendations = append(summary.Recommendations, "Check device and update health data")
	}

	// Анализируем каждую метрику
	summary.Temperature = analyzeTemperature(latestData.Temperature)
	summary.Sleep = analyzeSleep(latestData.SleepHours)
	summary.Activity = analyzeActivity(latestData.Activity)

	// Собираем проблемы из метрик
	if summary.Temperature.Status != "normal" {
		summary.Issues = append(summary.Issues, getTemperatureIssue(summary.Temperature))
		summary.Recommendations = append(summary.Recommendations, getTemperatureRecommendation(summary.Temperature))
	}

	if summary.Sleep.Status != "normal" {
		summary.Issues = append(summary.Issues, getSleepIssue(summary.Sleep))
		summary.Recommendations = append(summary.Recommendations, getSleepRecommendation(summary.Sleep))
	}

	if summary.Activity.Status != "normal" {
		summary.Issues = append(summary.Issues, getActivityIssue(summary.Activity))
		summary.Recommendations = append(summary.Recommendations, getActivityRecommendation(summary.Activity))
	}

	// Рассчитываем общий счет здоровья
	summary.HealthScore = calculateHealthScore(summary.Temperature.Score, summary.Sleep.Score, summary.Activity.Score)

	// Определяем общий статус и уровень предупреждения
	summary.OverallStatus, summary.AlertLevel, summary.NotificationLevel = determineOverallStatus(summary.HealthScore, len(summary.Issues))

	// Определяем необходимость внимания
	summary.RequiresAttention = summary.OverallStatus != "healthy" || summary.DataAge == "old"

	return summary
}

// analyzeTemperature анализирует температуру
func analyzeTemperature(temp float64) HealthMetric {
	metric := HealthMetric{
		Value:  temp,
		Normal: "37.5-39.5°C",
	}

	if temp >= NormalTemperatureMin && temp <= NormalTemperatureMax {
		metric.Status = "normal"
		metric.Score = 100
	} else if temp < NormalTemperatureMin-1.0 || temp > NormalTemperatureMax+1.0 {
		metric.Status = "critical"
		metric.Score = 20
	} else if temp < NormalTemperatureMin {
		metric.Status = "low"
		metric.Score = 60
	} else {
		metric.Status = "high"
		metric.Score = 60
	}

	return metric
}

// analyzeSleep анализирует сон
func analyzeSleep(sleep float64) HealthMetric {
	metric := HealthMetric{
		Value:  sleep,
		Normal: "8-16 hours",
	}

	if sleep >= NormalSleepMin && sleep <= NormalSleepMax {
		metric.Status = "normal"
		metric.Score = 100
	} else if sleep < NormalSleepMin-2.0 || sleep > NormalSleepMax+4.0 {
		metric.Status = "critical"
		metric.Score = 30
	} else if sleep < NormalSleepMin {
		metric.Status = "low"
		metric.Score = 70
	} else {
		metric.Status = "high"
		metric.Score = 70
	}

	return metric
}

// analyzeActivity анализирует активность
func analyzeActivity(activity float64) HealthMetric {
	metric := HealthMetric{
		Value:  activity,
		Normal: "30-80%",
	}

	// Нормальная активность 30-80%
	if activity >= 30 && activity <= 80 {
		metric.Status = "normal"
		metric.Score = 100
	} else if activity < 10 || activity > 95 {
		metric.Status = "critical"
		metric.Score = 25
	} else if activity < 30 {
		metric.Status = "low"
		metric.Score = 65
	} else {
		metric.Status = "high"
		metric.Score = 65
	}

	return metric
}

// calculateHealthScore рассчитывает общий балл здоровья
func calculateHealthScore(tempScore, sleepScore, activityScore int) int {
	// Взвешенное среднее: температура важнее всего
	return (tempScore*50 + sleepScore*30 + activityScore*20) / 100
}

// determineOverallStatus определяет общий статус здоровья
func determineOverallStatus(healthScore, issuesCount int) (status, alertLevel, notificationLevel string) {
	if healthScore >= 90 && issuesCount == 0 {
		return "healthy", "none", "info"
	} else if healthScore >= 70 && issuesCount <= 1 {
		return "minor_issues", "low", "info"
	} else if healthScore >= 50 {
		return "attention_needed", "medium", "warning"
	} else {
		return "critical", "high", "urgent"
	}
}

// Функции для получения описания проблем и рекомендаций
func getTemperatureIssue(metric HealthMetric) string {
	switch metric.Status {
	case "low":
		return "Body temperature is below normal"
	case "high":
		return "Body temperature is above normal"
	case "critical":
		return "Body temperature is critically abnormal"
	default:
		return "Temperature issue detected"
	}
}

func getTemperatureRecommendation(metric HealthMetric) string {
	switch metric.Status {
	case "low":
		return "Keep pet warm and consult veterinarian if temperature remains low"
	case "high":
		return "Ensure pet has access to cool areas and fresh water, monitor closely"
	case "critical":
		return "Contact veterinarian immediately - critical temperature detected"
	default:
		return "Monitor temperature closely"
	}
}

func getSleepIssue(metric HealthMetric) string {
	switch metric.Status {
	case "low":
		return "Pet is getting insufficient sleep"
	case "high":
		return "Pet is sleeping excessively"
	case "critical":
		return "Critical sleep pattern detected"
	default:
		return "Sleep pattern issue detected"
	}
}

func getSleepRecommendation(metric HealthMetric) string {
	switch metric.Status {
	case "low":
		return "Ensure quiet, comfortable sleeping environment and reduce stress"
	case "high":
		return "Increase activity and check for health issues causing lethargy"
	case "critical":
		return "Consult veterinarian about unusual sleep patterns"
	default:
		return "Monitor sleep patterns"
	}
}

func getActivityIssue(metric HealthMetric) string {
	switch metric.Status {
	case "low":
		return "Pet shows low activity levels"
	case "high":
		return "Pet shows unusually high activity"
	case "critical":
		return "Critical activity level detected"
	default:
		return "Activity level issue detected"
	}
}

func getActivityRecommendation(metric HealthMetric) string {
	switch metric.Status {
	case "low":
		return "Encourage gentle exercise and check for pain or illness"
	case "high":
		return "Monitor for stress or anxiety, ensure adequate rest periods"
	case "critical":
		return "Consult veterinarian about extreme activity changes"
	default:
		return "Monitor activity levels"
	}
}
