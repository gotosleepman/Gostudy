package models

import (
	"encoding/json"
	"fmt"

	"lending-copy/db"
)

type PoolBaseInfo struct {
	PoolID                 int             `json:"pool_id"`
	AutoLiquidateThreshold string          `json:"autoLiquidateThreshold"`
	BorrowSupply           string          `json:"borrowSupply"`
	BorrowToken            string          `json:"borrowToken"`
	BorrowTokenInfo        BorrowTokenInfo `json:"borrowTokenInfo"`
	EndTime                string          `json:"endTime"`
	InterestRate           string          `json:"interestRate"`
	JpCoin                 string          `json:"jpCoin"`
	LendSupply             string          `json:"lendSupply"`
	LendToken              string          `json:"lendToken"`
	LendTokenInfo          LendTokenInfo   `json:"lendTokenInfo"`
	MartgageRate           string          `json:"martgageRate"`
	MaxSupply              string          `json:"maxSupply"`
	SettleTime             string          `json:"settleTime"`
	SpCoin                 string          `json:"spCoin"`
	State                  string          `json:"state"`
}

type PoolBases struct {
	Id                     int    `json:"-" gorm:"column:id;primaryKey"`
	PoolID                 int    `json:"pool_id" gorm:"column:pool_id"`
	AutoLiquidateThreshold string `json:"autoLiquidateThreshold" gorm:"column:auto_liquidate_threshold"`
	BorrowSupply           string `json:"borrowSupply" gorm:"column:borrow_supply"`
	BorrowToken            string `json:"borrowToken" gorm:"column:borrow_token"`
	BorrowTokenInfo        string `json:"borrowTokenInfo" gorm:"column:borrow_token_info"`
	EndTime                string `json:"endTime" gorm:"column:end_time"`
	InterestRate           string `json:"interestRate" gorm:"column:interest_rate"`
	JpCoin                 string `json:"jpCoin" gorm:"column:jp_coin"`
	LendSupply             string `json:"lendSupply" gorm:"column:lend_supply"`
	LendToken              string `json:"lendToken" gorm:"column:lend_token"`
	LendTokenInfo          string `json:"lendTokenInfo" gorm:"column:lend_token_info"`
	MartgageRate           string `json:"martgageRate" gorm:"column:martgage_rate"`
	MaxSupply              string `json:"maxSupply" gorm:"column:max_supply"`
	SettleTime             string `json:"settleTime" gorm:"column:settle_time"`
	SpCoin                 string `json:"spCoin" gorm:"column:sp_coin"`
	State                  string `json:"state" gorm:"column:state"`
	ChainId                string `json:"chain_id" gorm:"column:chain_id"`
}

type BorrowTokenInfo struct {
	BorrowFee  string `json:"borrowFee"`
	TokenLogo  string `json:"tokenLogo"`
	TokenName  string `json:"tokenName"`
	TokenPrice string `json:"tokenPrice"`
}

type LendTokenInfo struct {
	LendFee    string `json:"lendFee"`
	TokenLogo  string `json:"tokenLogo"`
	TokenName  string `json:"tokenName"`
	TokenPrice string `json:"tokenPrice"`
}

type PoolBaseInfoRes struct {
	Index    int          `json:"index"`
	PoolData PoolBaseInfo `json:"pool_data"`
}

func NewPoolBases() *PoolBases {
	return &PoolBases{}
}

func (p *PoolBases) TableName() string {
	return "poolbases"
}

func (p *PoolBases) PoolBaseInfo(chainId int, res *[]PoolBaseInfoRes) error {
	var poolBases []PoolBases
	err := db.Mysql.Table("poolbases").Where("chain_id=?", fmt.Sprint(chainId)).Order("pool_id asc").Find(&poolBases).Error
	if err != nil {
		return err
	}
	for _, v := range poolBases {
		borrowTokenInfo := BorrowTokenInfo{}
		_ = json.Unmarshal([]byte(v.BorrowTokenInfo), &borrowTokenInfo)
		lendTokenInfo := LendTokenInfo{}
		_ = json.Unmarshal([]byte(v.LendTokenInfo), &lendTokenInfo)
		*res = append(*res, PoolBaseInfoRes{
			Index: v.PoolID - 1,
			PoolData: PoolBaseInfo{
				PoolID:                 v.PoolID,
				AutoLiquidateThreshold: v.AutoLiquidateThreshold,
				BorrowSupply:           v.BorrowSupply,
				BorrowToken:            v.BorrowToken,
				BorrowTokenInfo:        borrowTokenInfo,
				EndTime:                v.EndTime,
				InterestRate:           v.InterestRate,
				JpCoin:                 v.JpCoin,
				LendSupply:             v.LendSupply,
				LendToken:              v.LendToken,
				LendTokenInfo:          lendTokenInfo,
				MartgageRate:           v.MartgageRate,
				MaxSupply:              v.MaxSupply,
				SettleTime:             v.SettleTime,
				SpCoin:                 v.SpCoin,
				State:                  v.State,
			},
		})
	}
	return nil
}
