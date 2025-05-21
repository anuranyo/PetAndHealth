package data

type MasterDB interface {
	Users() UsersDB
	Pets() PetsDB
	Devices() DevicesDB
	Notifications() NotificationsDB
	HealthData() HealthDataDB
}
