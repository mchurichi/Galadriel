package spire

import (
	"context"

	"github.com/spiffe/go-spiffe/v2/bundle/spiffebundle"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	bundlev1 "github.com/spiffe/spire-api-sdk/proto/spire/api/server/bundle/v1"
	"github.com/spiffe/spire-api-sdk/proto/spire/api/types"
	"google.golang.org/grpc"
)

type fakeSpireBundleClient struct {
	bundle       *types.Bundle
	getBundleErr error
}

func (c fakeSpireBundleClient) CountBundles(ctx context.Context, in *bundlev1.CountBundlesRequest, opts ...grpc.CallOption) (*bundlev1.CountBundlesResponse, error) {
	return nil, nil
}

func (c fakeSpireBundleClient) GetBundle(ctx context.Context, in *bundlev1.GetBundleRequest, opts ...grpc.CallOption) (*types.Bundle, error) {
	if c.getBundleErr != nil {
		return nil, c.getBundleErr
	}

	return c.bundle, nil
}

func (fc fakeSpireBundleClient) AppendBundle(ctx context.Context, in *bundlev1.AppendBundleRequest, opts ...grpc.CallOption) (*types.Bundle, error) {
	return nil, nil
}

func (c fakeSpireBundleClient) PublishJWTAuthority(ctx context.Context, in *bundlev1.PublishJWTAuthorityRequest, opts ...grpc.CallOption) (*bundlev1.PublishJWTAuthorityResponse, error) {
	return nil, nil
}

func (c fakeSpireBundleClient) ListFederatedBundles(ctx context.Context, in *bundlev1.ListFederatedBundlesRequest, opts ...grpc.CallOption) (*bundlev1.ListFederatedBundlesResponse, error) {
	return nil, nil
}

func (c fakeSpireBundleClient) GetFederatedBundle(ctx context.Context, in *bundlev1.GetFederatedBundleRequest, opts ...grpc.CallOption) (*types.Bundle, error) {
	return nil, nil
}

func (c fakeSpireBundleClient) BatchCreateFederatedBundle(ctx context.Context, in *bundlev1.BatchCreateFederatedBundleRequest, opts ...grpc.CallOption) (*bundlev1.BatchCreateFederatedBundleResponse, error) {
	return nil, nil
}

func (c fakeSpireBundleClient) BatchUpdateFederatedBundle(ctx context.Context, in *bundlev1.BatchUpdateFederatedBundleRequest, opts ...grpc.CallOption) (*bundlev1.BatchUpdateFederatedBundleResponse, error) {
	return nil, nil
}
func (c fakeSpireBundleClient) BatchSetFederatedBundle(ctx context.Context, in *bundlev1.BatchSetFederatedBundleRequest, opts ...grpc.CallOption) (*bundlev1.BatchSetFederatedBundleResponse, error) {
	return nil, nil
}

func (c fakeSpireBundleClient) BatchDeleteFederatedBundle(ctx context.Context, in *bundlev1.BatchDeleteFederatedBundleRequest, opts ...grpc.CallOption) (*bundlev1.BatchDeleteFederatedBundleResponse, error) {
	return nil, nil
}

type fakeInternalClient struct {
	bundle                           *spiffebundle.Bundle
	federationRelationships          []*FederationRelationship
	getBundleErr                     error
	getFederationRelationshipsErr    error
	createFederationRelationshipsErr error
	updateFederationRelationshipsErr error
	deleteFederationRelationshipsErr error
}

func (c fakeInternalClient) GetBundle(context.Context) (*spiffebundle.Bundle, error) {
	if c.getBundleErr != nil {
		return nil, c.getBundleErr
	}

	return c.bundle, nil
}

func (c fakeInternalClient) ListFederationRelationships(context.Context) ([]*FederationRelationship, error) {
	if c.getFederationRelationshipsErr != nil {
		return nil, c.getFederationRelationshipsErr
	}

	return c.federationRelationships, nil
}

func (c fakeInternalClient) CreateFederationRelationships(context.Context, []*FederationRelationship) ([]*FederationRelationshipResult, error) {
	if c.createFederationRelationshipsErr != nil {
		return nil, c.createFederationRelationshipsErr
	}

	// TODO: add create logic and convert []*FederationRelationship to []*FederationRelationshipResult
	return []*FederationRelationshipResult{}, nil
}

func (c fakeInternalClient) UpdateFederationRelationships(context.Context, []*FederationRelationship) ([]*FederationRelationshipResult, error) {
	if c.updateFederationRelationshipsErr != nil {
		return nil, c.updateFederationRelationshipsErr
	}

	// TODO: add update logic
	return []*FederationRelationshipResult{}, nil
}

func (c fakeInternalClient) DeleteFederationRelationships(context.Context, []*spiffeid.TrustDomain) ([]*FederationRelationshipResult, error) {
	if c.deleteFederationRelationshipsErr != nil {
		return nil, c.deleteFederationRelationshipsErr
	}

	// TODO: add delete logic
	return []*FederationRelationshipResult{}, nil
}
