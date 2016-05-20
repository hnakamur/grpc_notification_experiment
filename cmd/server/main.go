package main

import (
	"flag"
	"io"
	"log"
	"net"

	"golang.org/x/net/context"

	pb "github.com/hnakamur/grpc_notification_experiment/sites"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

type sitesServer struct {
	modC chan *pb.SiteModification
	jobC chan *pb.Job
}

func newServer() *sitesServer {
	return &sitesServer{
		modC: make(chan *pb.SiteModification),
		jobC: make(chan *pb.Job),
	}
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

func (s *sitesServer) NotifySiteModification(ctx context.Context, mod *pb.SiteModification) (*pb.Empty, error) {
	s.modC <- mod
	return &pb.Empty{}, nil
}

func (s *sitesServer) WatchSites(_ *pb.Empty, stream pb.SitesService_WatchSitesServer) error {
	for mod := range s.modC {
		err := stream.Send(mod)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sitesServer) RequestWork(ctx context.Context, job *pb.Job) (*pb.Empty, error) {
	s.jobC <- job
	return &pb.Empty{}, nil
}

func (s *sitesServer) DoSomeWork(stream pb.SitesService_DoSomeWorkServer) error {
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("DoSomeWork handler got err=%s", err)
				break
			}

			log.Printf("res=%v", res)
		}
	}()
	for job := range s.jobC {
		err := stream.Send(job)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var enableTLS bool
	flag.BoolVar(&enableTLS, "enable-tls", false, "enable TLS")
	var certFile string
	flag.StringVar(&certFile, "cert-file", "../../ssl/server/server.crt", "TLS cert file")
	var keyFile string
	flag.StringVar(&keyFile, "key-file", "../../ssl/server/server.key", "TLS key file")
	var addr string
	flag.StringVar(&addr, "addr", ":10000", "server listen address")
	flag.Parse()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		grpclog.Fatal(err)
	}

	var opts []grpc.ServerOption
	if enableTLS {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSitesServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
