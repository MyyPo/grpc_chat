package main

import (
	"context"
	"log"
	"net"
	"os"
	"sync"

	chatpb "github.com/MyyPo/grpc-chat/pb/chat/v1"
	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
)

var grpcLog glog.LoggerV2

type Connection struct {
	stream chatpb.BroadcastService_ServerMessageStreamServer
	id     string
	active bool
	error  chan error
}

type Server struct {
	chatpb.UnimplementedBroadcastServiceServer
	// a collection of many client connections
	Connection []*Connection
}

// creates a new connection to the server and opens a channel sending errors
func (s *Server) ServerMessageStream(ptc *chatpb.ServerMessageStreamRequest, stream chatpb.BroadcastService_ServerMessageStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     ptc.User.Id,
		active: true,
		error:  make(chan error),
	}

	s.Connection = append(s.Connection, conn)

	return <-conn.error
}

func (s *Server) ClientMessage(ctx context.Context, msg *chatpb.Message) (*chatpb.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		// create a wait group based on the number of connections
		wait.Add(1)

		go func(msg *chatpb.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				grpcLog.Info("Sending message to: ", conn.stream)

				if err != nil {
					grpcLog.Errorf("Stream: %s stopped | Error: %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}
	// indentifies that the waitgroup has finished
	go func() {
		wait.Wait()
		close(done)
	}()
	// blocks until done channel gets closed
	<-done
	return &chatpb.Close{}, nil

}

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

func main() {
	var connections []*Connection

	server := &Server{chatpb.UnimplementedBroadcastServiceServer{}, connections}

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to create the server %v", err)
	}

	grpcLog.Info("The server successfuly started")

	chatpb.RegisterBroadcastServiceServer(grpcServer, server)
	grpcServer.Serve(lis)
}
