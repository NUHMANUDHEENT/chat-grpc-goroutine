package main

import (
    "log"
    "net"
    "sync"

    "google.golang.org/grpc"
    pb "chatservice/proto"
)

type chatServer struct {
    pb.UnimplementedChatServiceServer
    mu      sync.Mutex
    clients map[pb.ChatService_ChatServer]bool
}

func (s *chatServer) Chat(stream pb.ChatService_ChatServer) error {
    // Register new client
    s.mu.Lock()
    s.clients[stream] = true
    s.mu.Unlock()

    defer func() {
        s.mu.Lock()
        delete(s.clients, stream)
        s.mu.Unlock()
    }()

    // Listen for incoming messages from this client
    for {
        msg, err := stream.Recv()
        if err != nil {
            return err
        }

        // Broadcast the received message to all clients
        s.mu.Lock()
        for client := range s.clients {
            if client != stream {
                client.Send(msg)
            }
        }
        s.mu.Unlock()
    }
}

func main() {
    // Create a new gRPC server
    server := grpc.NewServer()
    chatService := &chatServer{clients: make(map[pb.ChatService_ChatServer]bool)}

    // Register the chat service
    pb.RegisterChatServiceServer(server, chatService)

    // Listen on port 50051
    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen on port 50051: %v", err)
    }

    log.Println("Server started on port :50051")
    if err := server.Serve(listener); err != nil {
        log.Fatalf("Failed to serve gRPC server: %v", err)
    }
}
