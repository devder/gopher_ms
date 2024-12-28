package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Data      string             `bson:"data" json:"data"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func New(mc *mongo.Client) Models {
	client = mc

	return Models{
		LogEntry: LogEntry{},
	}
}

func (l *LogEntry) Insert(entry LogEntry) error {
	coll := client.Database("logs").Collection("logs")

	_, err := coll.InsertOne(context.TODO(), entry)

	if err != nil {
		log.Println("Error inserting to into logs", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database("logs").Collection("logs")

	// sort by date
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "created_at", Value: -1}})
	cursor, err := coll.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	defer cursor.Close(ctx)
	var results []LogEntry
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Panic(err)
		return nil, err
	}

	return results, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry
	err = coll.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database("logs").Collection("logs")

	if err := coll.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(l.ID.Hex())
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": docID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: l.Name},
			{Key: "data", Value: l.Data},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	result, err := coll.UpdateOne(
		ctx,
		filter,
		update,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
