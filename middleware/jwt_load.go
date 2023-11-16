package middleware

import (
	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

func LoadDsrJwtClaims(ctx *gin.Context) (*DsrJwtClaims, error) {
	claims, ok := ctx.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		return nil, core.ErrClaimJWTFailed
	}

	customClaims, ok := claims.CustomClaims.(*DsrJwtClaims)
	if !ok {

		return nil, core.ErrClaimCustomJWTFailed
	}

	return customClaims, nil
}
