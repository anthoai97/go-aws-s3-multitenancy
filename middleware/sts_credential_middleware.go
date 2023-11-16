package middleware

import (
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
)

func (md *CustomMiddleware) CheckSTSCrendential(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		encounteredError := true

		AccessKey := c.Request.Header.Get("Access-Key-Id")
		SecretAccessKey := c.Request.Header.Get("Secret-Access-Key")
		SessionToken := c.Request.Header.Get("Session-Token")

		// We want to make sure the token is set, bail if not
		if AccessKey != "" && SecretAccessKey != "" && SessionToken != "" {
			logger.Debug("Authenticate STS Credential successful")
			encounteredError = false
		}

		if encounteredError {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				map[string]string{"message": "Missing sts credentials."},
			)
		}

		c.Set(CREDENTIAL_ACCESS_KEY, AccessKey)
		c.Set(CREDENTIAL_SECRET_KEY, SecretAccessKey)
		c.Set(CREDENTIAL_SESSION_KEY, SessionToken)

		// cred := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(AccessKey, SecretAccessKey, SessionToken))
		c.Next()
	}
}
