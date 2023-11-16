package api

import (
	"net/http"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/gin-gonic/gin"
)

func (api *api) UploadS3ObjectsByGenerateUrlHdl() func(*gin.Context) {
	return func(c *gin.Context) {

		// Parse JSON
		var json struct {
			Files []*entity.RequestFileUpload `json:"files" binding:"required"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()).WithDebug(err.Error()))
			return
		}

		if len(json.Files) < 1 {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(core.ErrBadRequest.ErrorField))
			return
		}

		res, err := api.business.UploadS3ObjectsByGenerateUrl(c, json.Files)
		if err != nil {
			core.WriteErrorResponse(c, err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
