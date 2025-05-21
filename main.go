package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/middle/tokenstore"

	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data/mongodb"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Помилка завантаження .env файлу")
	}

	ctx := context.Background()

	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	if username == "" || password == "" {
		panic("MONGO_USERNAME або MONGO_PASSWORD не встановлені")
	}

	// Формування URL для MongoDB
	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@petandhealth.dtpxu.mongodb.net/?retryWrites=true&w=majority&appName=PetAndHealth", username, password)
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Підключення до MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	petAndHealthDB := client.Database("PetAndHealth")
	mongoDB := mongodb.NewMasterDB(petAndHealthDB)

	// Start token cleanup routine
	startTokenCleanupRoutine()

	// Запуск API
	api.Run(api.Config{
		MasterDB: mongoDB,
	})
}

func startTokenCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			tokenstore.CleanupExpiredTokens()
		}
	}()
}
