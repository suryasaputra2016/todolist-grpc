package main

import (
	"context"
	"log"
	"net/http"

	"todolist-grpc/pb"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	e := echo.New()

	e.POST("/create-todo", func(c echo.Context) error {
		var req pb.CreateTodoRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		// Extract the JWT token from the Authorization header
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
		}

		// Create metadata with the token
		md := metadata.Pairs("authorization", token)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		// Connect to the gRPC server
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		defer conn.Close()

		// Create a new TodoService client
		client := pb.NewTodoServiceClient(conn)

		// Call the CreateTodo method
		res, err := client.CreateTodo(ctx, &req)
		if err != nil {
			log.Fatalf("Error while creating todo: %v", err)
		}

		// Return the response
		return c.JSON(http.StatusOK, res)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
