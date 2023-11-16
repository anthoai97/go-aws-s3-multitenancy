package api

import (
	"net/http"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/gin-gonic/gin"
)

func (api *api) ListS3StorageTreeHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		paramPairs := c.Request.URL.Query()
		path := paramPairs.Get("path")
		tenant := paramPairs.Get("tenant")
		nextContinuationToken := paramPairs.Get("next")

		if len(nextContinuationToken) > 0 {
			c.Set("NextContinuationToken", &nextContinuationToken)
		}

		if len(tenant) < 1 {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(core.ErrBadRequest.ErrorField).WithReason("tenant must be not empty"))
			return
		}

		tree, err := api.business.ListS3StorageTree(c, path, tenant)
		if err != nil {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		c.JSON(http.StatusOK, core.ResponseData(tree))
	}
}
