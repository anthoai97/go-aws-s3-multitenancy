package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WriteErrorResponse(c *gin.Context, err error) {
	if errSt, ok := err.(StatusCodeCarrier); ok {
		c.JSON(errSt.StatusCode(), errSt)
		return
	}

	c.JSON(http.StatusInternalServerError, ErrInternalServerError.WithError(err.Error()))
}
