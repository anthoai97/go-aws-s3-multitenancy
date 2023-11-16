package api

import (
	"net/http"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/gin-gonic/gin"
)

func (api *api) GenerateSTSCredentialHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		// Parse JSON
		var json entity.RequestSTSCredential

		if err := c.ShouldBindJSON(&json); err != nil {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(core.ErrBadRequest.ErrorField))
			return
		}

		if len(json.Tenant) < 1 {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(core.ErrBadRequest.ErrorField))
			return
		}

		cred, err := api.business.GenerateSTSCredential(c, json.Tenant)
		if err != nil {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		c.JSON(http.StatusOK, core.ResponseData(cred))
	}
}
