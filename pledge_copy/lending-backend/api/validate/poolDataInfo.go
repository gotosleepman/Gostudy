package validate

import (
	"io"

	"github.com/gin-gonic/gin"

	"lending-copy/api/common/statecode"
	"lending-copy/api/models/request"
)

type PoolDataInfo struct{}

func NewPoolDataInfo() *PoolDataInfo {
	return &PoolDataInfo{}
}

func (s *PoolDataInfo) PoolDataInfo(c *gin.Context, req *request.PoolDataInfo) int {
	err := c.ShouldBindQuery(req)
	if err == io.EOF {
		return statecode.ParameterEmptyErr
	} else if err != nil {
		return statecode.CommonErrServerErr
	}
	if req.ChainId != 97 && req.ChainId != 56 {
		return statecode.ChainIdErr
	}
	return statecode.CommonSuccess
}
