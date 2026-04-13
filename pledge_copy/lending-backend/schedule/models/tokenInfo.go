package models

import (
	"encoding/json"
	"errors"

	"lending-copy/db"

	"gorm.io/gorm"
)

func (TokenInfo) TableName() string { return "token_info" }

type TokenInfo struct {
	Id        int    `gorm:"column:id;primaryKey"`
	Logo      string `json:"logo" gorm:"column:logo"`
	Token     string `json:"token" gorm:"column:token"`
	Symbol    string `json:"symbol" gorm:"column:symbol"`
	ChainId   string `json:"chain_id" gorm:"column:chain_id"`
	Price     string `json:"price" gorm:"column:price"`
	Decimals  int    `json:"decimals" gorm:"column:decimals"`
	CreatedAt string `json:"created_at" gorm:"column:created_at"`
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at"`
}

func NewTokenInfo() *TokenInfo {
	return &TokenInfo{}
}

func (t *TokenInfo) GetTokenInfo(token, chainId string) (error, TokenInfo) {
	tokenInfo := TokenInfo{}
	redisKey := "token_info:" + chainId + ":" + token
	redisTokenInfoBytes, _ := db.RedisGet(redisKey)
	if len(redisTokenInfoBytes) <= 0 {
		err := db.Mysql.Table("token_info").Where("token=? and chain_id=?", token, chainId).First(&tokenInfo).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, tokenInfo
			}
			return err, tokenInfo
		}
		_ = db.RedisSet(redisKey, RedisTokenInfo{
			Token:   token,
			ChainId: chainId,
			Price:   tokenInfo.Price,
			Logo:    tokenInfo.Logo,
			Symbol:  tokenInfo.Symbol,
		}, 0)
		return nil, tokenInfo
	}
	redisTokenInfo := RedisTokenInfo{}
	err := json.Unmarshal(redisTokenInfoBytes, &redisTokenInfo)
	if err != nil {
		return err, tokenInfo
	}
	return nil, TokenInfo{
		Logo:    redisTokenInfo.Logo,
		Token:   token,
		Symbol:  redisTokenInfo.Symbol,
		ChainId: chainId,
		Price:   redisTokenInfo.Price,
	}
}
