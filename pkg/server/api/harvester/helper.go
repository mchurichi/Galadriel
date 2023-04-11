package harvester

import (
	"encoding/json"
	"io"

	"github.com/HewlettPackard/galadriel/pkg/common/entity"
	"github.com/labstack/echo/v4"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

func (b BundlePut) ToEntity() *entity.Bundle {
	return &entity.Bundle{
		Data:      []byte(b.TrustBundle),
		Signature: []byte(b.Signature),
		// TODO: do we need to store it in PEM or DER?
		SigningCertificate: []byte(b.SigningCertificate),
		TrustDomainName:    spiffeid.RequireTrustDomainFromString(b.TrustDomain),
	}
}

func (b *BundlePut) FromRequestBody(ctx echo.Context) error {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &b)
	if err != nil {
		return err
	}

	return nil
}
