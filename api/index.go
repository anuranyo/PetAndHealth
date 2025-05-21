package api

import (
	"context"
	"fmt"
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
	// Загрузка переменных окружения
	godotenv.Load()

	// Обработка путей - для Vercel важно корректно обрабатывать маршруты
	// Если запрос идет на /api/pet-and-health, удаляем этот префикс
	// Это позволит использовать те же маршруты, что и в локальной версии
	if strings.HasPrefix(r.URL.Path, "/api/pet-and-health") {
		r.URL.Path = r.URL.Path[18:] // Удаляем "/api/pet-and-health"
	} else if strings.HasPrefix(r.URL.Path, "/api") {
		r.URL.Path = r.URL.Path[4:] // Удаляем "/api"
	}

	// Если после удаления префикса путь пустой, установим его в "/"
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}

	// Создаем контекст для MongoDB
	ctx := context.Background()

	// Получаем учетные данные MongoDB из переменных окружения
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	if username == "" || password == "" {
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
		http.Error(w, "Ошибка подключения к MongoDB: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)

	// Проверяем соединение с MongoDB
	err = client.Ping(ctx, nil)
	if err != nil {
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
	router.ServeHTTP(w, r)
}
