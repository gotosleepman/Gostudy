package services

import (
	"lending-copy/api/common/statecode"
	"lending-copy/api/models"
	"lending-copy/api/models/request"
)

type TokenList struct{}

func NewTokenList() *TokenList {
	return &TokenList{}
}

func (c *TokenList) GetTokenList(req *request.TokenList) (int, []models.TokenList) {
	err, tokenList := models.NewTokenListModel().GetTokenList(req)
	if err != nil {
		return statecode.CommonErrServerErr, nil
	}
	return statecode.CommonSuccess, tokenList
}
