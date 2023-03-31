package test

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-storage/internal/grpc-storage/config"
	"grpc-storage/internal/grpc-storage/handler"
	"grpc-storage/internal/grpc-storage/storage"
	"grpc-storage/internal/grpc-storage/usecase"
	pb "grpc-storage/pkg/api/storage"
	"log"
	"net"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = config.New()
	logic := usecase.New(&storage.Storage{RAMStorage: make(map[string]string)})
	h := handler.New(logic)
	grpcServer := grpc.NewServer()

	log.Println("Starting Server ...")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:80"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	pb.RegisterStorageServer(grpcServer, h)

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("grpcServer Serve: %v", err)
		}

	}()

	m.Run()
	os.Exit(0)
}

func Test_SetGetValue(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	key := "eds"
	value := "qer"
	cl := pb.NewStorageClient(conn)

	set, err := cl.Set(context.Background(), &pb.SetRequest{Key: key, Value: value})
	if err != nil {
		t.Fatalf("%v", err)
	}

	if set.Status != "ok" {
		t.Fatal("Unexpected status code")
	}

	get, err := cl.Get(context.Background(), &pb.GetRequest{Key: key})
	if err != nil {
		t.Fatalf("%v", err)
	}

	if get.Value != value {
		t.Fatal("Wrong value")
	}
}
