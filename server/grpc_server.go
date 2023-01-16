package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"

	pb "weblibrary_sandbox/grpc_server"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"weblibrary_sandbox/database"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedGigaServiceServer
}

func (s *server) GetUser(ctx context.Context, id *pb.UserId) (*pb.User, error) {
	foundUser, err := database.GetUser(id.Id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return foundUser, nil
}

func (s *server) AddUser(ctx context.Context, newUser *pb.User) (*pb.UserId, error) {
	newUserID, err := database.AddUser(newUser.Name, int(newUser.Age))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return newUserID, nil
}

func (s *server) GetAllUsers(context.Context, *empty.Empty) (*pb.Users, error) {
	users, err := database.GetAllUsers()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return users, nil
}

func main() {

	err := database.InitDB()
	if err != nil {
		log.Fatalf("failed to establish connection to database: %v", err)
	}
	defer database.CloseDB()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGigaServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
