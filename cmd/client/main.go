package main

import (
	"flag"
	"io"
	"log"

	"golang.org/x/net/context"

	pb "bitbucket.org/hnakamur/grpc_notification_experiment/sites"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

func listSites(client pb.SitesServiceClient) {
	sites, err := client.ListSites(context.Background(), &pb.Empty{})
	if err != nil {
		grpclog.Fatalf("%v.ListSites(_) = _, %v: ", client, err)
	}
	for _, site := range sites.Sites {
		grpclog.Printf("site=%v", site)
	}
}

func watchSites(client pb.SitesServiceClient) {
	stream, err := client.WatchSites(context.Background(), &pb.Empty{})
	if err != nil {
		grpclog.Fatalf("%v.WatchSites(_) = _, %v: ", client, err)
	}
	for {
		mod, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			grpclog.Fatalf("%v.WatchSites(_) = _, %v: ", client, err)
		}

		grpclog.Printf("mod=%v", mod)
	}
}

func main() {
	var enableTLS bool
	flag.BoolVar(&enableTLS, "enable-tls", false, "enable TLS")
	var caFile string
	flag.StringVar(&caFile, "ca-file", "../../ssl/ca/cacert.pem", "The file containning the CA root cert file")
	var serverHostOverride string
	flag.StringVar(&serverHostOverride, "server-host-override", "grpc.example.com", "The server name use to verify the hostname returned by TLS handshake")
	var serverAddr string
	flag.StringVar(&serverAddr, "server-addr", "127.0.0.1:10000", "server listen address")
	flag.Parse()

	var opts []grpc.DialOption
	if enableTLS {
		var sn string
		if serverHostOverride != "" {
			sn = serverHostOverride
		}
		var creds credentials.TransportAuthenticator
		if caFile != "" {
			var err error
			creds, err = credentials.NewClientTLSFromFile(caFile, sn)
			if err != nil {
				grpclog.Fatalf("Failed to create TLS credentials %v", err)
			}
		} else {
			creds = credentials.NewClientTLSFromCert(nil, sn)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewSitesServiceClient(conn)
	listSites(client)
	watchSites(client)
}
