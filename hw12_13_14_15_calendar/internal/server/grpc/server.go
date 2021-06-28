package grpc

import (
	"net"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
)

// go generate: protoc -I="../../../api" go_out="../../../pkg/service" go-grpc_out="../../../pkg/service"

type grpcServer struct {
	log    logger.Logger
	server *grpc.Server
}

func NewGRPCServer(log logger.Logger) *grpcServer {
	return &grpcServer{log: log, server: grpc.NewServer()}
}

func (s *grpcServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.log.Info("grpc server started")

	return s.server.Serve(l)
}

func (s *grpcServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
		s.log.Info("grpc server stopped")
	}
}

func (s *grpcServer) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.server.RegisterService(desc, impl)
}
