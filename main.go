package main

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/spiffe/go-spiffe/v2/bundle/spiffebundle"
	"github.com/spiffe/spire-controller-manager/pkg/spireapi"
	"google.golang.org/grpc"
)

var spire1Socket = "/workspaces/Galadriel/dev/deployment/spire/spire1.sock"
var spire2Socket = "/workspaces/Galadriel/dev/deployment/spire/spire2.sock"

func main() {

	spire1 := newSpireServer(spire1Socket)
	spire2 := newSpireServer(spire2Socket)

	// SPIRE Server 1

	bundle1, err := spire1.GetBundle()
	if err != nil {
		panic(err)
	}
	pemBundle1, _ := bundle1.X509Bundle().Marshal()

	fmt.Println("Trust Domain 1:", bundle1.TrustDomain())
	fmt.Println("Bundle:")
	fmt.Println(string(pemBundle1))

	// SPIRE Server 2

	bundle2, err := spire2.GetBundle()
	if err != nil {
		panic(err)
	}
	pemBundle2, _ := bundle2.X509Bundle().Marshal()

	fmt.Println("Trust Domain 2:", bundle2.TrustDomain())
	fmt.Println("Bundle:")
	fmt.Println(string(pemBundle2))

	// Federate SPIRE 1 with SPIRE 2

	spire1.CreateFederationRelationship(bundle2)
}

type Client interface {
	spireapi.TrustDomainClient
	spireapi.BundleClient
}

type spireServer struct {
	client Client
}

func newSpireServer(socketPath string) *spireServer {
	client, err := dialSocket(context.Background(), socketPath)
	if err != nil {
		panic(err)
	}
	return &spireServer{
		client: client,
	}
}

func dialSocket(ctx context.Context, path string) (Client, error) {
	var target string
	if filepath.IsAbs(path) {
		target = "unix://" + path
	} else {
		target = "unix:" + path
	}
	grpcClient, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial API socket: %w", err)
	}

	return struct {
		spireapi.TrustDomainClient
		spireapi.BundleClient
		io.Closer
	}{
		TrustDomainClient: spireapi.NewTrustDomainClient(grpcClient),
		BundleClient:      spireapi.NewBundleClient(grpcClient),
		Closer:            grpcClient,
	}, nil
}

func (s *spireServer) GetBundle() (*spiffebundle.Bundle, error) {
	bundle, err := s.client.GetBundle(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to get bundle: %w", err)
	}
	return bundle, nil
}

func (s *spireServer) GetFederationRelationships() ([]spireapi.FederationRelationship, error) {
	feds, err := s.client.ListFederationRelationships(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to list federation relationships: %w", err)
	}
	return feds, nil
}

func (s *spireServer) CreateFederationRelationship(bundle *spiffebundle.Bundle) (*spireapi.Status, error) {
	x509bundle := bundle.X509Bundle()

	fmt.Println("Creating federation relationship with", bundle.TrustDomain().ID())
	spireSpiffeId, _ := bundle.TrustDomain().ID().AppendPath("/spire/server")

	status, err := s.client.CreateFederationRelationships(context.TODO(), []spireapi.FederationRelationship{
		{
			TrustDomain:       x509bundle.TrustDomain(),
			TrustDomainBundle: bundle,
			BundleEndpointURL: "https://localhost:8442",
			BundleEndpointProfile: spireapi.HTTPSSPIFFEProfile{
				EndpointSPIFFEID: spireSpiffeId,
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create federation relationship: %w", err)
	}

	return &status[0], nil // why many?
}
