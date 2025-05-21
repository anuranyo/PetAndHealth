package mongodb

import (
	"context"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const NotificationsCollectionName = "Notifications"

type notificationsDB struct {
	collection *mongo.Collection
}

func newNotificationsDB(db *mongo.Database) *notificationsDB {
	return &notificationsDB{
		collection: db.Collection(NotificationsCollectionName),
	}
}

func NewNotificationsDB(db *mongo.Database) data.NotificationsDB {
	return newNotificationsDB(db)
}

func (n *notificationsDB) Get(id string) (*data.Notification, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result data.Notification
	err = n.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (n *notificationsDB) Insert(notif *data.Notification) error {
	_, err := n.collection.InsertOne(context.TODO(), notif)
	return err
}

func (n *notificationsDB) Update(id string, updateFields bson.M) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = n.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.M{"$set": updateFields},
	)
	return err
}

func (n *notificationsDB) GetForUser(userID primitive.ObjectID) ([]*data.Notification, error) {
	cursor, err := n.collection.Find(
		context.TODO(),
		bson.M{"user_id": userID},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var notifications []*data.Notification
	if err = cursor.All(context.TODO(), &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (n *notificationsDB) MarkAsRead(id primitive.ObjectID) error {
	_, err := n.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"read": true}},
	)
	if err != nil {
		return err
	}
	return nil
}

func (n *notificationsDB) GetUnreadCount(userID primitive.ObjectID) (int, error) {
	count, err := n.collection.CountDocuments(
		context.TODO(),
		bson.M{
			"user_id": userID,
			"read":    false,
		},
	)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
