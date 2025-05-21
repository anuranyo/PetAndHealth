package api

import (
	"net/http"
	"os"

	"context"
	"fmt"

	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data/mongodb"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Handler обрабатывает запросы Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Загрузка переменных окружения
	godotenv.Load()

	ctx := context.Background()

	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	if username == "" || password == "" {
		http.Error(w, "MONGO_USERNAME или MONGO_PASSWORD не установлены", http.StatusInternalServerError)
		return
	}

	// Формирование URL для MongoDB
	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@petandhealth.dtpxu.mongodb.net/?retryWrites=true&w=majority&appName=PetAndHealth", username, password)
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Подключение к MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)

	petAndHealthDB := client.Database("PetAndHealth")
	mongoDB := mongodb.NewMasterDB(petAndHealthDB)

	// Используем GetRouter вместо Run
	router := api.GetRouter(api.Config{
		MasterDB: mongoDB,
	})

	// Обрабатываем текущий запрос
	router.ServeHTTP(w, r)
}
