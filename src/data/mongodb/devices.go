package mongodb

import (
	"context"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const DevicesCollectionName = "Devices"

type devicesDB struct {
	collection *mongo.Collection
}

func newDevicesDB(db *mongo.Database) *devicesDB {
	return &devicesDB{
		collection: db.Collection(DevicesCollectionName),
	}
}

func NewDevicesDB(db *mongo.Database) data.DevicesDB {
	return newDevicesDB(db)
}

func (d *devicesDB) Get(id primitive.ObjectID) (*data.Device, error) {
	var result data.Device
	err := d.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *devicesDB) Insert(device *data.Device) error {
	if device.LastSyncTime.T == 0 {
		epochSeconds := time.Now().Unix()
		device.LastSyncTime = primitive.Timestamp{
			T: uint32(epochSeconds),
			I: 0,
		}
	}

	_, err := d.collection.InsertOne(context.TODO(), device)
	return err
}

func (d *devicesDB) Update(id primitive.ObjectID, updateFields bson.M) error {
	if _, ok := updateFields["last_sync_time"]; ok {
		if epochSeconds, ok := updateFields["last_sync_time"].(int64); ok {
			updateFields["last_sync_time"] = time.Unix(epochSeconds, 0)
		}
	}

	_, err := d.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": updateFields},
	)
	return err
}

func (d *devicesDB) GetAll() ([]*data.Device, error) {
	var devices []*data.Device
	cursor, err := d.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var device *data.Device
		err := cursor.Decode(&device)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}
