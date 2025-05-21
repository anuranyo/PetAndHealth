package requests

import (
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func GetUserName(db data.MasterDB, ownerID primitive.ObjectID) string {
	usersDB := db.Users()
	owner, err := usersDB.Get(ownerID)
	if err != nil {
		log.Printf("Failed to retrieve owner with ID %s: %v", ownerID.Hex(), err)
		return "Unknown"
	}
	return owner.FullName
}
