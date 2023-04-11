package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/HewlettPackard/galadriel/pkg/common/api"
	"github.com/HewlettPackard/galadriel/pkg/common/entity"
	harvesterAPI "github.com/HewlettPackard/galadriel/pkg/server/api/harvester"
	"github.com/HewlettPackard/galadriel/pkg/server/datastore"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const tokenKey = "token"

type HarvesterAPIHandlers struct {
	Logger    logrus.FieldLogger
	Datastore datastore.Datastore
}

func NewHarvesterAPIHandlers(l logrus.FieldLogger, ds datastore.Datastore) HarvesterAPIHandlers {
	return HarvesterAPIHandlers{
		Logger:    l,
		Datastore: ds,
	}
}

func (h *HarvesterAPIHandlers) GetRelationships(ctx echo.Context, params harvesterAPI.GetRelationshipsParams) error {
	return nil
}

func (h *HarvesterAPIHandlers) PatchRelationshipsRelationshipID(ctx echo.Context, relationshipID uuid.UUID) error {
	return nil
}

func (h *HarvesterAPIHandlers) Onboard(ctx echo.Context) error {
	return nil
}

func (h HarvesterAPIHandlers) BundleSync(ctx echo.Context, trustDomainName api.TrustDomainName) error {
	return nil
}

func (h *HarvesterAPIHandlers) BundlePut(ctx echo.Context, trustDomainName api.TrustDomainName) error {
	h.Logger.Debug("Receiving post bundle request")

	// TODO: move authn out and replace with Access Token when implemented
	jt, ok := ctx.Get(tokenKey).(*entity.JoinToken)
	if !ok {
		err := errors.New("error parsing token")
		h.handleTCPError(ctx, err.Error())
		return err
	}

	token, err := h.Datastore.FindJoinToken(ctx.Request().Context(), jt.Token)
	if err != nil {
		err := errors.New("error looking up token")
		h.handleTCPError(ctx, err.Error())
		return err
	}

	authenticatedTD, err := h.Datastore.FindTrustDomainByID(ctx.Request().Context(), token.TrustDomainID)
	if err != nil {
		err := errors.New("error looking up trust domain")
		h.handleTCPError(ctx, err.Error())
		return err
	}

	if authenticatedTD.Name.String() != trustDomainName {
		return fmt.Errorf("authenticated trust domain {%s} does not match trust domain in path: {%s}", authenticatedTD.Name, trustDomainName)
	}
	// end authn

	req := &harvesterAPI.BundlePut{}
	req.FromRequestBody(ctx)

	if authenticatedTD.Name.String() != req.TrustDomain {
		err := fmt.Errorf("authenticated trust domain {%s} does not match trust domain in request body: {%s}", authenticatedTD.Name, req.TrustDomain)
		h.handleTCPError(ctx, err.Error())
		return err
	}

	storedBundle, err := h.Datastore.FindBundleByTrustDomainID(ctx.Request().Context(), authenticatedTD.ID.UUID)
	if err != nil {
		h.handleTCPError(ctx, err.Error())
		return err
	}

	if req.TrustBundle == "" {
		return nil
	}

	bundle := req.ToEntity()
	if storedBundle != nil {
		bundle.TrustDomainID = storedBundle.TrustDomainID
	}
	res, err := h.Datastore.CreateOrUpdateBundle(ctx.Request().Context(), bundle)
	if err != nil {
		h.handleTCPError(ctx, err.Error())
		return err
	}

	if err = WriteResponse(ctx, res); err != nil {
		h.handleTCPError(ctx, err.Error())
		return err
	}

	return nil
}

func (h *HarvesterAPIHandlers) handleTCPError(ctx echo.Context, errMsg string) {
	h.Logger.Errorf(errMsg)
	_, err := ctx.Response().Write([]byte(errMsg))
	if err != nil {
		h.Logger.Errorf("Failed to write error response: %v", err)
	}
}

// Writes the response to the client
func WriteResponse(ctx echo.Context, body interface{}) error {
	if body == nil {
		return nil
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal response body: %v", err)
	}
	_, err = ctx.Response().Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write response body: %v", err)
	}

	return nil
}

// func (e *Endpoints) syncFederatedBundleHandler(ctx echo.Context) error {
// 	e.Logger.Debug("Receiving sync request")

// 	jt, ok := ctx.Get(tokenKey).(*entity.JoinToken)
// 	if !ok {
// 		err := errors.New("error parsing join token")
// 		e.handleTCPError(ctx, err.Error())
// 		return err
// 	}

// 	token, err := e.Datastore.FindJoinToken(ctx.Request().Context(), jt.Token)
// 	if err != nil {
// 		err := errors.New("error looking up token")
// 		e.handleTCPError(ctx, err.Error())
// 		return err
// 	}

// 	harvesterTrustDomain, err := e.Datastore.FindTrustDomainByID(ctx.Request().Context(), token.TrustDomainID)
// 	if err != nil {
// 		e.handleTCPError(ctx, fmt.Sprintf("failed to read body: %v", err))
// 		return err
// 	}

// 	body, err := io.ReadAll(ctx.Request().Body)
// 	if err != nil {
// 		e.handleTCPError(ctx, fmt.Sprintf("failed to read body: %v", err))
// 		return err
// 	}

// 	receivedHarvesterState := common.SyncBundleRequest{}
// 	err = json.Unmarshal(body, &receivedHarvesterState)
// 	if err != nil {
// 		e.handleTCPError(ctx, fmt.Sprintf("failed to unmarshal state: %v", err))
// 		return err
// 	}

// 	harvesterBundleDigests := receivedHarvesterState.State

// 	_, foundSelf := receivedHarvesterState.State[harvesterTrustDomain.Name]
// 	if foundSelf {
// 		e.handleTCPError(ctx, "bad request: harvester cannot federate with itself")
// 		return err
// 	}

// 	relationships, err := e.Datastore.FindRelationshipsByTrustDomainID(ctx.Request().Context(), harvesterTrustDomain.ID.UUID)
// 	if err != nil {
// 		e.handleTCPError(ctx, fmt.Sprintf("failed to fetch relationships: %v", err))
// 		return err
// 	}

// 	federatedTDs := getFederatedTrustDomains(relationships, harvesterTrustDomain.ID.UUID)

// 	if len(federatedTDs) == 0 {
// 		e.Logger.Debug("No federated trust domains yet")
// 		return nil
// 	}

// 	federatedBundles, federatedBundlesDigests, err := e.getCurrentFederatedBundles(ctx.Request().Context(), federatedTDs)
// 	if err != nil {
// 		e.handleTCPError(ctx, fmt.Sprintf("failed to fetch bundles from DB: %v", err))
// 		return err
// 	}

// 	if len(federatedBundles) == 0 {
// 		e.Logger.Debug("No federated bundles yet")
// 		return nil
// 	}

// 	bundlesUpdates, err := e.getFederatedBundlesUpdates(ctx.Request().Context(), harvesterBundleDigests, federatedBundles)
// 	if err != nil {
// 		e.handleTCPError(ctx, fmt.Sprintf("failed to fetch bundles from DB: %v", err))
// 		return err
// 	}

// 	response := common.SyncBundleResponse{
// 		Updates: bundlesUpdates,
// 		State:   federatedBundlesDigests,
// 	}

// 	responseBytes, err := json.Marshal(response)
// 	if err != nil {
// 		e.handleTCPError(ctx, fmt.Sprintf("failed to marshal response: %v", err))
// 		return err
// 	}

// 	_, err = ctx.Response().Write(responseBytes)
// 	if err != nil {
// 		e.handleTCPError(ctx, fmt.Sprintf("failed to write response: %v", err))
// 		return err
// 	}

// 	return nil
// }

// func getFederatedTrustDomains(relationships []*entity.Relationship, tdID uuid.UUID) []uuid.UUID {
// 	var federatedTrustDomains []uuid.UUID

// 	for _, r := range relationships {
// 		ma := r.TrustDomainAID
// 		mb := r.TrustDomainBID

// 		if tdID == ma {
// 			federatedTrustDomains = append(federatedTrustDomains, mb)
// 		} else {
// 			federatedTrustDomains = append(federatedTrustDomains, ma)
// 		}
// 	}
// 	return federatedTrustDomains
// }

// func (e *Endpoints) getFederatedBundlesUpdates(ctx context.Context, harvesterBundlesDigests common.BundlesDigests, federatedBundles []*entity.Bundle) (common.BundleUpdates, error) {
// 	response := make(common.BundleUpdates)

// 	for _, b := range federatedBundles {
// 		td, err := e.Datastore.FindTrustDomainByID(ctx, b.TrustDomainID)
// 		if err != nil {
// 			return nil, err
// 		}

// 		serverDigest := util.GetDigest(b.Data)
// 		harvesterDigest := harvesterBundlesDigests[td.Name]

// 		// If the bundle digest received from a federated trust domain of the calling harvester is not the same as the
// 		// digest the server has, the harvester needs to be updated of the new bundle. This also covers the case of
// 		// the harvester not being aware of any bundles. The update represents a newly federated trustDomain's bundle.
// 		if !bytes.Equal(harvesterDigest, serverDigest) {
// 			response[td.Name] = b
// 		}
// 	}

// 	return response, nil
// }

// func (e *Endpoints) getCurrentFederatedBundles(ctx context.Context, federatedTDs []uuid.UUID) ([]*entity.Bundle, common.BundlesDigests, error) {
// 	var bundles []*entity.Bundle
// 	bundlesDigests := map[spiffeid.TrustDomain][]byte{}

// 	for _, id := range federatedTDs {
// 		b, err := e.Datastore.FindBundleByTrustDomainID(ctx, id)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		td, err := e.Datastore.FindTrustDomainByID(ctx, id)
// 		if err != nil {
// 			return nil, nil, err
// 		}

// 		if b != nil {
// 			bundles = append(bundles, b)
// 			bundlesDigests[td.Name] = util.GetDigest(b.Data)
// 		}
// 	}

// 	return bundles, bundlesDigests, nil
// }
