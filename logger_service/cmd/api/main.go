package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/devder/gopher_ms/logger/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongodb:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	log.Printf("Starting logger service on port %s\n", webPort)
	// connect to mongo
	var err error
	client, err = connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	// create a ctx in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// start rpc server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Panic(err)
	}
	go app.rpcListen()
	// start web server
	app.serve()
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to mongo:", err)
		return nil, err
	}

	log.Println("Connected to MongoDB")

	return c, nil
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) rpcListen() {
	log.Printf("Starting RPC server on %s", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Panicf("Error starting rpc server: %s", err)
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			log.Println("Error accepting connection")
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
