package mongodb

import (
	"context"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const PetsCollectionName = "Pets"

type petsDB struct {
	collection *mongo.Collection
}

func newPetsDB(db *mongo.Database) *petsDB {
	return &petsDB{
		collection: db.Collection(PetsCollectionName),
	}
}

func NewPetsDB(db *mongo.Database) data.PetsDB {
	return newPetsDB(db)
}

func (p *petsDB) createPetAggregationPipeline(matchStage bson.D) mongo.Pipeline {
	return mongo.Pipeline{
		{
			{"$match", matchStage},
		},
		{
			{"$lookup", bson.D{
				{"from", UsersCollectionName},
				{"localField", "owner_id"},
				{"foreignField", "_id"},
				{"as", "owner"},
			}},
		},
		{
			{"$lookup", bson.D{
				{"from", HealthDataCollectionName},
				{"localField", "_id"},
				{"foreignField", "pet_id"},
				{"as", "health_data"},
			}},
		},
		{
			{"$unwind", bson.D{
				{"path", "$owner"},
				{"preserveNullAndEmptyArrays", true},
			}},
		},
		{
			{"$group", bson.D{
				{"_id", "$_id"},
				{"name", bson.D{{"$first", "$name"}}},
				{"species", bson.D{{"$first", "$species"}}},
				{"breed", bson.D{{"$first", "$breed"}}},
				{"age", bson.D{{"$first", "$age"}}},
				{"owner_id", bson.D{{"$first", "$owner_id"}}},
				{"owner", bson.D{{"$first", "$owner"}}},
				{"health_data", bson.D{{"$first", "$health_data"}}},
			}},
		},
	}
}

func (p *petsDB) executeAggregation(pipeline mongo.Pipeline) ([]*data.Pet, error) {
	cursor, err := p.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var pets []*data.Pet
	if err := cursor.All(context.TODO(), &pets); err != nil {
		return nil, err
	}

	return pets, nil
}

func (p *petsDB) Get(id primitive.ObjectID) (*data.Pet, error) {
	pipeline := p.createPetAggregationPipeline(bson.D{{"_id", id}})

	pets, err := p.executeAggregation(pipeline)
	if err != nil {
		return nil, err
	}

	if len(pets) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return pets[0], nil
}

func (p *petsDB) Insert(pet *data.Pet) error {
	petToInsert := &data.Pet{
		ID:      pet.ID,
		Name:    pet.Name,
		Species: pet.Species,
		Breed:   pet.Breed,
		Age:     pet.Age,
		OwnerID: pet.OwnerID,
	}

	_, err := p.collection.InsertOne(context.TODO(), petToInsert)
	return err
}

func (p *petsDB) Update(id primitive.ObjectID, updateFields bson.M) error {
	delete(updateFields, "health")
	delete(updateFields, "owner")

	_, err := p.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": updateFields},
	)
	return err
}

func (p *petsDB) Delete(id primitive.ObjectID) error {
	_, err := p.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func (p *petsDB) GetAll() ([]*data.Pet, error) {
	pipeline := p.createPetAggregationPipeline(bson.D{})
	return p.executeAggregation(pipeline)
}

func (p *petsDB) GetByFilter(filter bson.M) ([]*data.Pet, error) {
	matchStage := bson.D{}
	for key, value := range filter {
		matchStage = append(matchStage, bson.E{Key: key, Value: value})
	}

	pipeline := p.createPetAggregationPipeline(matchStage)
	return p.executeAggregation(pipeline)
}
