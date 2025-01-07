// gRPC Client (client.go)
package main

import (
    "context"
    "fmt"
    "google.golang.org/grpc"
    "rssreader"
    "time"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        fmt.Printf("Failed to connect: %v", err)
        return
    }
    defer conn.Close()

    client := rssreader.NewRssServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    urls := []string{"https://example.com/rss"}
    req := &rssreader.ParseRequest{Urls: urls}
    res, err := client.ParseUrls(ctx, req)
    if err != nil {
        fmt.Printf("Error calling gRPC service: %v", err)
        return
    }

    fmt.Printf("Received response: %v\n", res.Items)
}
