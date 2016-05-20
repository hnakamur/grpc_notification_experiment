package main

import (
	"flag"
	"io"
	"log"
	"strings"

	"golang.org/x/net/context"

	pb "github.com/hnakamur/grpc_notification_experiment/sites"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

type sitesClient struct {
	client   pb.SitesServiceClient
	clientID string
}

func newSitesClient(client pb.SitesServiceClient, clientID string) *sitesClient {
	return &sitesClient{
		client:   client,
		clientID: clientID,
	}
}

func (c *sitesClient) listSites(ctx context.Context) {
	sites, err := c.client.ListSites(context.Background(), &pb.Empty{})
	if err != nil {
		grpclog.Fatalf("%v.ListSites(_) = _, %v: ", c.client, err)
	}
	for _, site := range sites.Sites {
		grpclog.Printf("site=%v", site)
	}
}

func (c *sitesClient) notifySiteModification(ctx context.Context, op, domain, origin string) {
	var opVal pb.SiteModificationOp
	switch op {
	case "add":
		opVal = pb.SiteModificationOp_ADDED
	case "edit":
		opVal = pb.SiteModificationOp_EDITED
	case "remove":
		opVal = pb.SiteModificationOp_REMOVED
	default:
		log.Fatal("operation must be one of add, edit or remove")
	}
	mod := &pb.SiteModification{
		Op: opVal,
		Site: &pb.Site{
			Domain: domain,
			Origin: origin,
		},
	}
	_, err := c.client.NotifySiteModification(context.Background(), mod)
	if err != nil {
		grpclog.Fatalf("%v.NotifySiteModification(_) = _, %v: ", c.client, err)
	}
}

func (c *sitesClient) watchSites(ctx context.Context) {
	stream, err := c.client.WatchSites(ctx, &pb.Empty{})
	if err != nil {
		grpclog.Fatalf("%v.WatchSites(_) = _, %v: ", c.client, err)
	}
	for {
		mod, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			grpclog.Fatalf("%v.WatchSites(_) = _, %v: ", c.client, err)
		}

		grpclog.Printf("mod=%v", mod)
	}
}

func (c *sitesClient) requestWork(ctx context.Context, targets []string) {
	job := &pb.Job{Targets: targets}
	_, err := c.client.RequestWork(ctx, job)
	if err != nil {
		grpclog.Fatalf("%v.RequestWork(_) = _, %v: ", c.client, err)
	}
}

func (c *sitesClient) doSomeWork(ctx context.Context) {
	stream, err := c.client.DoSomeWork(ctx)
	if err != nil {
		grpclog.Fatalf("%v.DoSomeWork(_) = _, %v: ", c.client, err)
	}
	for {
		job, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			grpclog.Fatalf("%v.DoSomeWork(_) = _, %v: ", c.client, err)
		}

		grpclog.Printf("job=%v", job)

		results := make([]*pb.JobResult, 0, len(job.Targets))
		for _, target := range job.Targets {
			results = append(results,
				&pb.JobResult{
					Target: target,
					Result: "success",
				})
		}
		result := &pb.JobResults{
			ClientID: c.clientID,
			Results:  results,
		}
		err = stream.Send(result)
		if err != nil {
			grpclog.Fatalf("%v.DoSomeWork(_) = _, %v: ", c.client, err)
		}
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
	var op string
	flag.StringVar(&op, "op", "watch", "operation: one of work, req, watch, add, remove, edit")
	var domain string
	flag.StringVar(&domain, "domain", "example.com", "domain of site")
	var origin string
	flag.StringVar(&origin, "origin", "example.org", "origin of site")
	var clientID string
	flag.StringVar(&clientID, "client-id", "client1", "client ID")
	var targets string
	flag.StringVar(&targets, "targets", "target1,target2", "comma seperated targets")
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
	client := newSitesClient(pb.NewSitesServiceClient(conn), clientID)
	ctx := context.Background()
	switch op {
	case "watch":
		client.listSites(ctx)
		client.watchSites(ctx)
	case "req":
		client.requestWork(ctx, strings.Split(targets, ","))
	case "work":
		client.doSomeWork(ctx)
	default:
		client.notifySiteModification(ctx, op, domain, origin)
	}
}
