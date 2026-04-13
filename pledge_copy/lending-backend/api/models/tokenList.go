package models

import (
	"errors"
	"fmt"

	"lending-copy/api/models/request"
	"lending-copy/db"
)

type TokenList struct {
	Id       int32  `json:"-" gorm:"column:id;primaryKey"`
	Symbol   string `json:"symbol" gorm:"column:symbol"`
	Decimals int    `json:"decimals" gorm:"column:decimals"`
	Token    string `json:"token" gorm:"column:token"`
	Logo     string `json:"logo" gorm:"column:logo"`
	ChainId  string `json:"chain_id" gorm:"column:chain_id"`
}

func (TokenList) TableName() string { return "token_info" }

func NewTokenListModel() *TokenListModel {
	return &TokenListModel{}
}

type TokenListModel struct{}

func (m *TokenListModel) GetTokenList(req *request.TokenList) (error, []TokenList) {
	var tokenList []TokenList
	cid := fmt.Sprint(req.ChainId)
	err := db.Mysql.Table("token_info").Where("chain_id = ?", cid).Find(&tokenList).Error
	if err != nil {
		return errors.New("record select err " + err.Error()), nil
	}
	return nil, tokenList
}
