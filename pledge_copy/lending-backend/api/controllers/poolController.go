package controllers

import (
	"regexp"
	"strings"
	"time"

	"lending-copy/api/common/statecode"
	"lending-copy/api/models"
	"lending-copy/api/models/request"
	"lending-copy/api/models/response"
	"lending-copy/api/services"
	"lending-copy/api/validate"
	"lending-copy/config"

	"github.com/gin-gonic/gin"
)

type PoolController struct{}

func (c *PoolController) PoolBaseInfo(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.PoolBaseInfo{}
	var result []models.PoolBaseInfoRes
	errCode := validate.NewPoolBaseInfo().PoolBaseInfo(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	errCode = services.NewPool().PoolBaseInfo(req.ChainId, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	res.Response(ctx, statecode.CommonSuccess, result)
}

func (c *PoolController) PoolDataInfo(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.PoolDataInfo{}
	var result []models.PoolDataInfoRes
	errCode := validate.NewPoolDataInfo().PoolDataInfo(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	errCode = services.NewPool().PoolDataInfo(req.ChainId, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	res.Response(ctx, statecode.CommonSuccess, result)
}

func (c *PoolController) TokenList(ctx *gin.Context) {
	req := request.TokenList{}
	result := response.TokenList{}
	errCode := validate.NewTokenList().TokenList(ctx, &req)
	if errCode != statecode.CommonSuccess {
		ctx.JSON(200, map[string]string{"error": "chainId error"})
		return
	}
	errCode, data := services.NewTokenList().GetTokenList(&req)
	if errCode != statecode.CommonSuccess {
		ctx.JSON(200, map[string]string{"error": "chainId error"})
		return
	}
	baseURL := c.GetBaseURL()
	result.Name = "Lending Copy Token List"
	result.LogoURI = baseURL + "storage/img/logo.png"
	result.Timestamp = time.Now()
	result.Version = response.Version{Major: 1, Minor: 0, Patch: 0}
	for _, v := range data {
		result.Tokens = append(result.Tokens, response.Token{
			Name:     v.Symbol,
			Symbol:   v.Symbol,
			Decimals: v.Decimals,
			Address:  v.Token,
			ChainID:  v.ChainId,
			LogoURI:  v.Logo,
		})
	}
	ctx.JSON(200, result)
}

func (c *PoolController) Search(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.Search{}
	result := response.Search{}
	errCode := validate.NewSearch().Search(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	errCode, count, pools := services.NewSearch().Search(&req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	result.Rows = pools
	result.Count = count
	res.Response(ctx, statecode.CommonSuccess, result)
}

func (c *PoolController) GetBaseURL() string {
	domainName := config.Config.Env.DomainName
	if domainName == "" {
		return config.Config.Env.Protocol + "://127.0.0.1:" + config.Config.Env.Port + "/"
	}
	domainNameSlice := strings.Split(domainName, "")
	pattern := `^\d+`
	isNumber, _ := regexp.MatchString(pattern, domainNameSlice[0])
	if isNumber {
		return config.Config.Env.Protocol + "://" + config.Config.Env.DomainName + ":" + config.Config.Env.Port + "/"
	}
	return config.Config.Env.Protocol + "://" + config.Config.Env.DomainName + "/"
}
