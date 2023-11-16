package middleware

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

// checkJWT is a gin.HandlerFunc middleware
// that will check the validity of our JWT.
func CheckJWT() gin.HandlerFunc {
	authDomain := core.GetEnvVar("AUTH0_DOMAIN", "dataspire-dev.eu.auth0.com")
	audience := core.GetEnvVar("AUTH0_AUDIENCE", "https://api.datastorage.dataspire.com/")
	issuerURL, err := url.Parse("https://" + authDomain + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	// Set up the validator.
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &DsrJwtClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)

	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(ctx *gin.Context) {
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			ctx.Request = r
			ctx.Next()
		}

		middleware.CheckJWT(handler).ServeHTTP(ctx.Writer, ctx.Request)

		if encounteredError {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				map[string]string{"message": "JWT is invalid."},
			)
		}
	}
}
