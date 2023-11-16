package middleware

import (
	"github.com/ethereum/go-ethereum/log"
)

var CREDENTIAL_ACCESS_KEY = "AccessKey"
var CREDENTIAL_SECRET_KEY = "SecretKey"
var CREDENTIAL_SESSION_KEY = "SessionToken"

type CustomMiddleware struct {
	Logger log.Logger
}
