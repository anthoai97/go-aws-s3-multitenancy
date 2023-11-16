package middleware

import (
	"context"
	"errors"
)

// DsrJwtClaims contains custom data we want from the token.
type DsrJwtClaims struct {
	OrganisationID string `json:"organisation_id"`
}

// Validate errors out if `OrganisationID` is empty.
func (c *DsrJwtClaims) Validate(ctx context.Context) error {
	if len(c.OrganisationID) < 1 {
		return errors.New("should reject was set to true")
	}
	return nil
}
