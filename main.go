package main

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/spiffe/go-spiffe/v2/bundle/spiffebundle"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/spire-controller-manager/pkg/spireapi"
	"google.golang.org/grpc"
)

var spire1Socket = "/Users/mchurichi/Documents/scytale/src/spire/api1.sock"
var spire2Socket = "/Users/mchurichi/Documents/scytale/src/spire/api2.sock"

func main() {
	client1, err := dialSocket(context.Background(), spire1Socket)
	if err != nil {
		panic(err)
	}
	client2, err := dialSocket(context.Background(), spire2Socket)
	if err != nil {
		panic(err)
	}

	// Get Trust Domain Bundle from SPIRE 1
	bundle1, err := getBundle(context.Background(), client1)
	if err != nil {
		panic(err)
	}
	x509bundle1 := bundle1.X509Bundle()
	pemBundle1, err := x509bundle1.Marshal()
	if err != nil {
		panic(err)
	}
	fmt.Println("Bundle:")
	fmt.Println("Trust Domain:", x509bundle1.TrustDomain())
	fmt.Println(string(pemBundle1))

	// Get Trust Domain Bundle from SPIRE 2
	bundle2, err := getBundle(context.Background(), client2)
	if err != nil {
		panic(err)
	}
	x509bundle2 := bundle2.X509Bundle()
	pemBundle2, err := x509bundle2.Marshal()
	if err != nil {
		panic(err)
	}
	fmt.Println("Bundle:")
	fmt.Println("Trust Domain:", x509bundle2.TrustDomain())
	fmt.Println(string(pemBundle2))

	// List federation relationships
	feds, err := listFederationRelationships(context.Background(), client1)
	if err != nil {
		panic(err)
	}
	fmt.Println("Federation Relationships:")
	fmt.Println(feds)

	// Create SPIRE 1 federation relationship
	spire2SpiffeID, err := spiffeid.FromString("spiffe://two.org/spire/server")
	if err != nil {
		panic(err)
	}

	fmt.Println("Creating federation relationship with spiffe://two.org/spire/server")
	status1, err := client1.CreateFederationRelationships(context.Background(), []spireapi.FederationRelationship{
		{
			TrustDomain:       x509bundle2.TrustDomain(),
			TrustDomainBundle: bundle2,
			BundleEndpointURL: "https://localhost:8442",
			BundleEndpointProfile: spireapi.HTTPSSPIFFEProfile{
				EndpointSPIFFEID: spire2SpiffeID,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Creation status:", status1)
	fmt.Println(status1)

	// List federation relationships
	feds, err = listFederationRelationships(context.Background(), client1)
	if err != nil {
		panic(err)
	}
	fmt.Println("Federation Relationships:")
	fmt.Println(feds)
}

func listFederationRelationships(ctx context.Context, client Client) ([]spireapi.FederationRelationship, error) {
	feds, err := client.ListFederationRelationships(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list federation relationships: %w", err)
	}
	return feds, nil
}

func getBundle(ctx context.Context, client Client) (*spiffebundle.Bundle, error) {
	bundle, err := client.GetBundle(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get bundle: %w", err)
	}
	return bundle, nil
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

type Client interface {
	spireapi.TrustDomainClient
	spireapi.BundleClient
}
