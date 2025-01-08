package main

import (
	"context"
	"log"
	"net"

	"your_project_path/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type server struct {
	pb.UnimplementedTodoServiceServer
	todos map[string]*pb.Todo
}

func (s *server) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	// Extract the token from the metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing metadata")
	}

	token := md["authorization"]
	if len(token) == 0 {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing token")
	}

	// Validate the token (implement your own validation logic)
	if !validateToken(token[0]) {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}

	// Create the todo item
	id := generateID() // Implement a function to generate unique IDs
	todo := &pb.Todo{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}
	s.todos[id] = todo

	return &pb.CreateTodoResponse{Id: id}, nil
}

func validateToken(token string) bool {
	// Implement your token validation logic here
	// For example, you can use a JWT library to validate the token
	return true
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, &server{todos: make(map[string]*pb.Todo)})

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
