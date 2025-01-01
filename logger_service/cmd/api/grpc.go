package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/devder/gopher_ms/logger/data"
	"github.com/devder/gopher_ms/logger/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	res := &logs.LogResponse{Result: "logged"}
	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()
	logs.RegisterLogServiceServer(grpcServer, &LogServer{Models: app.Models})
	// consider this self documentation to help the client
	reflection.Register(grpcServer)

	log.Printf("start gRPC Server at %s", lis.Addr().String())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
