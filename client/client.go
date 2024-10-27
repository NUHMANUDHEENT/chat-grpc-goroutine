package main

import (
    "bufio"
    "context"
    "log"
    "os"
    "time"

    "google.golang.org/grpc"
    pb "chatservice/proto"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("Failed to connect to server: %v", err)
    }
    defer conn.Close()

    client := pb.NewChatServiceClient(conn)
    stream, err := client.Chat(context.Background())
    if err != nil {
        log.Fatalf("Error creating stream: %v", err)
    }

    go func() {
        for {
            msg, err := stream.Recv()
            if err != nil {
                log.Fatalf("Error receiving message: %v", err)
            }
            log.Printf("%s: %s", msg.User, msg.Message)
        }
    }()

    reader := bufio.NewReader(os.Stdin)
    for {
        log.Print("Enter your name: ")
        user, _ := reader.ReadString('\n')
        user = user[:len(user)-1]

        for {
            log.Print("Enter your message: ")
            text, _ := reader.ReadString('\n')
            text = text[:len(text)-1]

            msg := &pb.ChatMessage{
                User:    user,
                Message: text,
            }

            if err := stream.Send(msg); err != nil {
                log.Fatalf("Error sending message: %v", err)
            }

            time.Sleep(time.Second)
        }
    }
}
