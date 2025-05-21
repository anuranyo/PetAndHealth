package handlers

import (
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/requests"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func GetPetReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Отримуємо ID тварини з URL параметрів
	petIDStr := chi.URLParam(r, "id")
	if petIDStr == "" {
		http.Error(w, "Pet ID is required", http.StatusBadRequest)
		return
	}

	// Конвертуємо ID в ObjectID
	petID, err := primitive.ObjectIDFromHex(petIDStr)
	if err != nil {
		http.Error(w, "Invalid pet ID format", http.StatusBadRequest)
		return
	}

	// Парсимо тіло запиту для отримання часових рамок звіту
	req, err := requests.NewPetReportTimeRange(r)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.StartTime == 0 || req.EndTime == 0 {
		http.Error(w, "Missing start_time or end_time parameters", http.StatusBadRequest)
		return
	}

	startTime := primitive.Timestamp{T: uint32(req.StartTime), I: 0}
	endTime := primitive.Timestamp{T: uint32(req.EndTime), I: 0}

	petsDB := MongoDB(r).Pets()
	healthDB := MongoDB(r).HealthData()

	// Отримуємо інформацію про тварину
	pet, err := petsDB.Get(petID)
	if err != nil {
		http.Error(w, "Failed to retrieve pet information", http.StatusInternalServerError)
		return
	}

	// Формуємо фільтр для отримання даних про здоров'я
	filter := bson.M{
		"pet_id": petID,
		"time": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	// Отримуємо дані про здоров'я
	healthData, err := healthDB.GetByFilter(filter)
	if err != nil {
		http.Error(w, "Failed to retrieve health data", http.StatusInternalServerError)
		return
	}

	// Отримуємо ім'я власника
	ownerName := requests.GetUserName(MongoDB(r), pet.OwnerID)

	// Генеруємо PDF звіт
	pdf, err := GeneratePetReportPDF(pet, healthData, ownerName)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	// Встановлюємо заголовки відповіді
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="pet_report.pdf"`)
	w.WriteHeader(http.StatusOK)

	// Відправляємо PDF як відповідь
	err = pdf.Output(w)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}
}
