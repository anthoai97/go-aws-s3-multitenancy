package core

import (
	"net/http"

	logger "github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
)

var Log = logger.New("logscope", "error")

func WriteErrorResponse(c *gin.Context, err error) {
	if errSt, ok := err.(StatusCodeCarrier); ok {
		c.JSON(errSt.StatusCode(), errSt)
		return
	}

	Log.Error("WriteErrorResponse", "Error", err)
	c.JSON(http.StatusInternalServerError, ErrInternalServerError.WithError(err.Error()))
}
