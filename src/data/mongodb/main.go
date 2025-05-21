package mongodb

import (
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"go.mongodb.org/mongo-driver/mongo"
)

type masterDB struct {
	users         *usersDB
	pets          *petsDB
	devices       *devicesDB
	notifications *notificationsDB
	healthData    *healthDataDB
}

func NewMasterDB(db *mongo.Database) data.MasterDB {
	return &masterDB{
		users:         newUsersDB(db),
		pets:          newPetsDB(db),
		devices:       newDevicesDB(db),
		notifications: newNotificationsDB(db),
		healthData:    newHealthDataDB(db),
	}
}

func (m *masterDB) Users() data.UsersDB {
	return m.users
}

func (m *masterDB) Pets() data.PetsDB {
	return m.pets
}

func (m *masterDB) Devices() data.DevicesDB {
	return m.devices
}

func (m *masterDB) Notifications() data.NotificationsDB {
	return m.notifications
}

func (m *masterDB) HealthData() data.HealthDataDB {
	return m.healthData
}
