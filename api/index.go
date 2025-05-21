package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data/mongodb"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Handler - точка входа для Vercel бессерверных функций
func Handler(w http.ResponseWriter, r *http.Request) {
	// Логирование для отладки
	log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

	// Загрузка переменных окружения
	godotenv.Load()

	// ВАЖНО: Корректируем путь для Vercel
	originalPath := r.URL.Path

	// Удаляем "/api" в начале пути - это критически важно для Vercel
	if strings.HasPrefix(originalPath, "/api") {
		r.URL.Path = strings.TrimPrefix(originalPath, "/api")
	}

	// Если путь пустой, установим "/"
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}

	// Добавляем префикс "/api/pet-and-health" для соответствия вашей маршрутизации
	r.URL.Path = "/api/pet-and-health" + r.URL.Path

	log.Printf("Измененный путь: %s", r.URL.Path)

	// Создаем контекст для MongoDB
	ctx := context.Background()

	// Получаем учетные данные MongoDB из переменных окружения
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	if username == "" || password == "" {
		log.Printf("Ошибка: отсутствуют учетные данные MongoDB")
		http.Error(w, "MONGO_USERNAME или MONGO_PASSWORD не установлены", http.StatusInternalServerError)
		return
	}

	// Формирование URL для MongoDB
	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@petandhealth.dtpxu.mongodb.net/?retryWrites=true&w=majority&appName=PetAndHealth", username, password)
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Подключение к MongoDB с таймаутом для бессерверных функций
	clientOptions.SetConnectTimeout(10 * time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Ошибка подключения к MongoDB: %v", err)
		http.Error(w, "Ошибка подключения к MongoDB: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)

	// Проверяем соединение с MongoDB
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Не удалось подключиться к MongoDB: %v", err)
		http.Error(w, "Не удалось установить соединение с MongoDB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Инициализируем базу данных
	petAndHealthDB := client.Database("PetAndHealth")
	mongoDB := mongodb.NewMasterDB(petAndHealthDB)

	// Получаем маршрутизатор с настроенными маршрутами
	router := api.GetRouter(api.Config{
		MasterDB: mongoDB,
	})

	// Обрабатываем запрос с помощью маршрутизатора
	log.Printf("Передаем запрос маршрутизатору: %s %s", r.Method, r.URL.Path)
	router.ServeHTTP(w, r)
}
