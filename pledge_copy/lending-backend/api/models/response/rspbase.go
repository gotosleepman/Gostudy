package response

import (
	"lending-copy/api/common/statecode"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	Res *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

type Page struct {
	Code  int         `json:"code"`
	Msg   string      `json:"message"`
	Total int         `json:"total"`
	Data  interface{} `json:"data"`
}

func (g *Gin) Response(c *gin.Context, code int, data interface{}, httpStatus ...int) {
	lang := statecode.LangEn
	if langInf, ok := c.Get("lang"); ok {
		lang = langInf.(int)
	}
	rsp := Response{
		Code: code,
		Msg:  statecode.GetMsg(code, lang),
		Data: data,
	}
	status := 200
	if len(httpStatus) > 0 {
		status = httpStatus[0]
	}
	g.Res.JSON(status, rsp)
}

func (g *Gin) ResponsePages(c *gin.Context, code int, totalCount int, data interface{}) {
	lang := statecode.LangEn
	if langInf, ok := c.Get("lang"); ok {
		lang = langInf.(int)
	}
	g.Res.JSON(200, Page{
		Code:  code,
		Msg:   statecode.GetMsg(code, lang),
		Total: totalCount,
		Data:  data,
	})
}
