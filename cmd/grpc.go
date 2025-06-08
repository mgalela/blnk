package main

import (
	"context"
	"net"

	"github.com/jerry-enebeli/blnk"
	"github.com/jerry-enebeli/blnk/blnkgrpc"
	"github.com/jerry-enebeli/blnk/config"
	pb "github.com/jerry-enebeli/blnk/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// var kasp = keepalive.ServerParameters{
// 	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
// 	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
// 	MaxConnectionAgeGrace: 10 * time.Second, // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
// 	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
// 	Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
// }

// var kaep = keepalive.EnforcementPolicy{
// 	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
// 	PermitWithoutStream: true,            // Allow pings even when there are no active streams
// }

func startGrpcServer(ctx context.Context, b *blnk.Blnk, cfg config.ServerConfig) error {
	lis, err := net.Listen("tcp", ":" + cfg.GrpcPort)
	if err != nil {
		logrus.Warningf("Failed GRPC to listen: %v", err.Error())
		return nil
	}
	s := grpc.NewServer()
	// s := grpc.NewServer(
	// 		grpc.KeepaliveEnforcementPolicy(kaep),
	// 		grpc.KeepaliveParams(kasp),
	// 		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	// 	)

	bgrpc := blnkgrpc.NewGRPC(ctx, b)
	pb.RegisterBlnkServiceServer(s, bgrpc)

	// Register reflection service on gRPC server.
	 reflection.Register(s)
	 if err := s.Serve(lis); err != nil {
	 	logrus.Warningf("Failed to serve: %v", err.Error())
		return nil
	 } 
 
	return nil

}

