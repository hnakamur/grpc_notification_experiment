package main

import (
	"flag"
	"net"
	"time"

	"golang.org/x/net/context"

	pb "bitbucket.org/hnakamur/grpc_notification_experiment/sites"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type sitesServer struct{}

func newServer() *sitesServer {
	return &sitesServer{}
}

func (s *sitesServer) ListSites(ctx context.Context, _ *pb.Empty) (*pb.Sites, error) {
	sites := &pb.Sites{
		Sites: []*pb.Site{
			&pb.Site{Domain: "foo.example.com", Origin: "foo.example.org"},
			&pb.Site{Domain: "bar.example.com", Origin: "bar.example.org"},
		},
	}
	return sites, nil
}

func (s *sitesServer) WatchSites(_ *pb.Empty, stream pb.SitesService_WatchSitesServer) error {
	mod := &pb.SiteModification{
		Op:   pb.SiteModificationOp_EDITED,
		Site: &pb.Site{Domain: "foo.example.com", Origin: "foo.example.net"},
	}
	err := stream.Send(mod)
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	mod = &pb.SiteModification{
		Op:   pb.SiteModificationOp_ADDED,
		Site: &pb.Site{Domain: "baz.example.com", Origin: "baz.example.org"},
	}
	err = stream.Send(mod)
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	mod = &pb.SiteModification{
		Op:   pb.SiteModificationOp_REMOVED,
		Site: &pb.Site{Domain: "baz.example.com", Origin: "baz.example.org"},
	}
	err = stream.Send(mod)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":10000", "server listen address")
	flag.Parse()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		grpclog.Fatal(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSitesServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
