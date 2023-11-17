package api

import (
	"net/http"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/gin-gonic/gin"
)

func (api *api) DownloadS3ObjectsByGenerateUrlHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		// Parse JSON
		var json *entity.RequestFileDownload

		if err := c.ShouldBindJSON(&json); err != nil {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()).WithDebug(err.Error()))
			return
		}

		if len(json.FilePath) < 1 || len(json.Tenant) < 1 {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(core.ErrBadRequest.ErrorField))
			return
		}

		url, err := api.business.DownloadS3ObjectsByGenerateUrl(c, json)
		if err != nil {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()).WithDebug(err.Error()))
			return
		}

		// Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
		c.JSON(http.StatusOK, core.ResponseData(url))

	}
}
