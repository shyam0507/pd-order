package storage

import (
	"context"
	"time"

	"github.com/shyam0507/pd-order/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client *mongo.Client
	dbName string
}

const Collection_Name = "order"

// CreateOrder implements Storage.
func (m MongoStorage) CreateOrder(o types.Order) error {

	m.client.Database(m.dbName).Collection(Collection_Name).InsertOne(
		context.Background(),
		o,
	)

	return nil
}

// UpdateOrder implements Storage.
func (m MongoStorage) UpdateOrder(id string, status string) error {

	oId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}
	_, err = m.client.Database(m.dbName).Collection(Collection_Name).UpdateOne(
		context.Background(),
		bson.M{"_id": oId},
		bson.M{"$set": bson.M{"status": status}},
	)

	return err
}

func NewMongoStorage(uri string, dbName string) Storage {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		return nil
	}

	return MongoStorage{client: client, dbName: dbName}
}
