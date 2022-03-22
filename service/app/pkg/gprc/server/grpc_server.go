package grpcserver

import (
	"fmt"
	"github.com/vladazn/dhq/proto/gen/go/proto/storage"
	"github.com/vladazn/dhq/service/app/service"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	host     string
	port     int
	server   *grpc.Server
	services *service.Services
}

func NewGrpcServer(
	host string,
	port int,
	services *service.Services,
) Server {

	return Server{
		host:     host,
		port:     port,
		server:   grpc.NewServer(),
		services: services,
	}
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		return err
	}

	storage.RegisterStorageServer(s.server, newStorageServer(s.services))

	err = s.server.Serve(lis)

	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
