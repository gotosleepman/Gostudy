package validate

import (
	"io"

	"github.com/gin-gonic/gin"

	"lending-copy/api/common/statecode"
	"lending-copy/api/models/request"
)

type PoolBaseInfo struct{}

func NewPoolBaseInfo() *PoolBaseInfo {
	return &PoolBaseInfo{}
}

func (s *PoolBaseInfo) PoolBaseInfo(c *gin.Context, req *request.PoolBaseInfo) int {
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
