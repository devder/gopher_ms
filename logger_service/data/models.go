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
