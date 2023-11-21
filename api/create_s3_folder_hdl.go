package api

import (
	"net/http"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/gin-gonic/gin"
)

func (api *api) CreateS3FolderHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		// Parse JSON
		var requestBody entity.RequestCreateFolder

		if err := c.ShouldBind(&requestBody); err != nil {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()).WithDebug(err.Error()))
			return
		}

		if len(requestBody.Tenant) < 1 || len(requestBody.Path) < 1 {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(core.ErrBadRequest.ErrorField))
			return
		}

		res, err := api.business.CreateS3Folder(c, &requestBody)
		if err != nil {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()).WithDebug(err.Error()))
			return
		}

		c.JSON(http.StatusOK, core.ResponseData(res))
	}
}
