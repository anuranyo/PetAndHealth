package utils

import (
	"fmt"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

const (
	MinTemperature = 37.0
	MaxTemperature = 39.5
	MinSleepHours  = 7.0
	MaxSleepHours  = 14.0
)

func CheckPetHealthAndNotify(notificationsDB data.NotificationsDB, ownerID primitive.ObjectID,
	temperature float64, sleepHours float64, currentTime time.Time) {

	if temperature < MinTemperature || temperature > MaxTemperature {
		message := fmt.Sprintf("Abnormal temperature detected: %.2f°C. Please check your pet's health.", temperature)
		title := "Temperature Warning"
		CreateHealthNotification(notificationsDB, ownerID, title, message, "temperature_warning", currentTime)
	}

	if sleepHours > MaxSleepHours || sleepHours < MinSleepHours {
		message := fmt.Sprintf("Abnormal sleep hours detected: %.2f hours. Please monitor your pet's activity.", sleepHours)
		title := "Sleep Pattern Warning"
		CreateHealthNotification(notificationsDB, ownerID, title, message, "sleep_warning", currentTime)
	}
}

func CreateHealthNotification(notificationsDB data.NotificationsDB, ownerID primitive.ObjectID,
	title string, message string, notifType string, currentTime time.Time) {

	notification := data.Notification{
		ID:        primitive.NewObjectID(),
		UserID:    ownerID,
		Title:     title,
		Message:   message,
		Type:      notifType,
		Time:      currentTime,
		Read:      false,
		Delivered: false,
	}

	err := notificationsDB.Insert(&notification)
	if err != nil {
		log.Printf("Failed to insert notification: %s\n", err)
	}
}
