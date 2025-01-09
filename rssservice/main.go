// Main file to run a basic gRPC server (main.go)

package main

import (
    "fmt"
    "net"
    "google.golang.org/grpc"
)

func main() {
    lis, err := net.Listen("tcp", ":50051")

    if err != nil {
        fmt.Printf("Failed to listen: %v", err)
        return
    }
    
	fmt.Println("gRPC server running on port 50051")
    
	grpcServer := grpc.NewServer()

	// Register RssService here

    if err := grpcServer.Serve(lis); err != nil {
        fmt.Printf("Failed to serve: %v", err)
    }
}