package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/devder/gopher_ms/logger/data"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	coll := client.Database("logs").Collection("logs")
	_, err := coll.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("error writing to mongo :", err)
		return err
	}

	*resp = fmt.Sprintf("Processed payload via RPC: %s", payload.Name)

	return nil
}
