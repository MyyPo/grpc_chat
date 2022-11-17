package service

import (
	"context"
	chatpb "github.com/MyyPo/grpc-chat/pb/chat/v1"
	glog "google.golang.org/grpc/grpclog"
	"sync"
)

var connections []*Connection

func NewChatServer(grpcLog glog.LoggerV2) *ChatServer {
	return &ChatServer{
		chatpb.UnimplementedBroadcastServiceServer{},
		connections,
		grpcLog,
	}
}

type Connection struct {
	stream chatpb.BroadcastService_ServerMessageStreamServer
	id     string
	active bool
	error  chan error
}

type ChatServer struct {
	chatpb.UnimplementedBroadcastServiceServer
	// a collection of many client connections
	Connection []*Connection
	grpcLog    glog.LoggerV2
}

// creates a new connection to the server and opens a channel sending errors
func (s *ChatServer) ServerMessageStream(ptc *chatpb.ServerMessageStreamRequest, stream chatpb.BroadcastService_ServerMessageStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     ptc.User.Id,
		active: true,
		error:  make(chan error),
	}

	s.Connection = append(s.Connection, conn)

	return <-conn.error
}

func (s *ChatServer) ClientMessage(ctx context.Context, msg *chatpb.Message) (*chatpb.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		// create a wait group based on the number of connections
		wait.Add(1)

		go func(msg *chatpb.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				s.grpcLog.Info("Sending message to: ", conn.stream)

				if err != nil {
					s.grpcLog.Errorf("Stream: %s stopped | Error: %v", conn.stream, err)
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
