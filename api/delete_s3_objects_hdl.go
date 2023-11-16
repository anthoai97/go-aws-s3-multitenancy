package api

import (
	"errors"
	"net/http"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/gin-gonic/gin"
)

func (api *api) DeleteS3ObjectsHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		// Parse JSON
		var json struct {
			Objects []*entity.RequestObjectDelete `json:"objects" binding:"required"`
		}

		if err := c.ShouldBind(&json); err != nil {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()).WithDebug(err.Error()))
			return
		}

		if len(json.Objects) < 1 {
			err := errors.New("paths is emmpty")
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()).WithDebug(err.Error()))
			return
		}

		res, err := api.business.DeleteS3Objects(c, json.Objects)
		if err != nil {
			core.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()).WithDebug(err.Error()))
			return
		}

		c.JSON(http.StatusOK, core.ResponseData(res))
	}
}
