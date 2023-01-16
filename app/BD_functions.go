package main

import (
	"context"
	"flag"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	pb "weblibrary_sandbox/grpc_server"
	"weblibrary_sandbox/models"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func convertUser(user *pb.User) models.User {
	return models.User{
		Name:   user.GetName(),
		Age:    int(user.GetAge()),
		UserID: uuid.MustParse(user.UserId),
	}
}

func getAllUsers() ([]models.User, error) {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGigaServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	// add intentional delay to show timeout error handling
	if rand.Intn(100) > 80 {
		time.Sleep(time.Millisecond * 220)
	}

	allUsersResponse, err := c.GetAllUsers(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	var allUsers []models.User

	for _, v := range allUsersResponse.GetUsers() {
		allUsers = append(allUsers, convertUser(v))
	}

	return allUsers, nil
}

func getUser(id uuid.UUID) (models.User, error) {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGigaServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	// add intentional delay to show timeout error handling
	if rand.Intn(100) > 80 {
		time.Sleep(time.Millisecond * 220)
	}

	userID := &pb.UserId{Id: id.String()}
	foundUser, err := c.GetUser(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	return convertUser(foundUser), nil
}

func addUser(newUser models.User) error {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err)
		return err
	}
	defer conn.Close()
	c := pb.NewGigaServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	user := &pb.User{Name: newUser.Name, Age: int32(newUser.Age)}
	_, err = c.AddUser(ctx, user)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
