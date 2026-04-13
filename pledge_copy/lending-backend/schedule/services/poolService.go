package services

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"

	"lending-copy/config"
	"lending-copy/contract/bindings"
	"lending-copy/db"
	"lending-copy/log"
	"lending-copy/schedule/models"
	"lending-copy/utils"
)

type poolService struct{}

func NewPool() *poolService {
	return &poolService{}
}

func (s *poolService) UpdateAllPoolInfo() {
	s.UpdatePoolInfo(config.Config.TestNet.LendingPoolAddr, config.Config.TestNet.NetUrl, config.Config.TestNet.ChainId)
}

func (s *poolService) UpdatePoolInfo(contractAddress, network, chainId string) {
	if contractAddress == "" || contractAddress == "0x0000000000000000000000000000000000000000" {
		log.Logger.Sugar().Warn("UpdatePoolInfo skipped: lending_pool_addr not configured")
		return
	}
	log.Logger.Sugar().Info("UpdatePoolInfo ", contractAddress, network)
	cli, err := bindings.Dial(network, contractAddress)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	defer cli.Close()

	ctx := context.Background()
	borrowFee, err := cli.BorrowFee(ctx)
	if err != nil {
		log.Logger.Sugar().Error("BorrowFee ", err)
		return
	}
	lendFee, err := cli.LendFee(ctx)
	if err != nil {
		log.Logger.Sugar().Error("LendFee ", err)
		return
	}
	pLength, err := cli.PoolLength(ctx)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	n := int(pLength.Int64())
	for i := 0; i < n; i++ {
		poolId := utils.IntToString(i + 1)
		baseInfo, err := cli.PoolBaseInfo(ctx, big.NewInt(int64(i)))
		if err != nil {
			log.Logger.Sugar().Info("UpdatePoolInfo PoolBaseInfo err", poolId, err)
			continue
		}
		_, borrowToken := models.NewTokenInfo().GetTokenInfo(baseInfo.BorrowToken.Hex(), chainId)
		_, lendToken := models.NewTokenInfo().GetTokenInfo(baseInfo.LendToken.Hex(), chainId)
		lendTokenJSON, _ := json.Marshal(models.LendToken{
			LendFee:    lendFee.String(),
			TokenLogo:  lendToken.Logo,
			TokenName:  lendToken.Symbol,
			TokenPrice: lendToken.Price,
		})
		borrowTokenJSON, _ := json.Marshal(models.BorrowToken{
			BorrowFee:  borrowFee.String(),
			TokenLogo:  borrowToken.Logo,
			TokenName:  borrowToken.Symbol,
			TokenPrice: borrowToken.Price,
		})
		poolBase := models.PoolBase{
			SettleTime:             baseInfo.SettleTime.String(),
			PoolId:                 utils.StringToInt(poolId),
			ChainId:                chainId,
			EndTime:                baseInfo.EndTime.String(),
			InterestRate:           baseInfo.InterestRate.String(),
			MaxSupply:              baseInfo.MaxSupply.String(),
			LendSupply:             baseInfo.LendSupply.String(),
			BorrowSupply:           baseInfo.BorrowSupply.String(),
			MartgageRate:           baseInfo.MartgageRate.String(),
			LendToken:              baseInfo.LendToken.Hex(),
			LendTokenInfo:          string(lendTokenJSON),
			BorrowToken:            baseInfo.BorrowToken.Hex(),
			BorrowTokenInfo:        string(borrowTokenJSON),
			State:                  utils.IntToString(int(baseInfo.State)),
			SpCoin:                 baseInfo.SpCoin.Hex(),
			JpCoin:                 baseInfo.JpCoin.Hex(),
			AutoLiquidateThreshold: baseInfo.AutoLiquidateThreshold.String(),
		}
		hasInfoData, byteBaseInfoStr, baseInfoMd5Str := s.GetPoolMd5(&poolBase, "base_info:lc_pool_"+chainId+"_"+poolId)
		if !hasInfoData || baseInfoMd5Str != byteBaseInfoStr {
			err = models.NewPoolBase().SavePoolBase(chainId, poolId, &poolBase)
			if err != nil {
				log.Logger.Sugar().Error("SavePoolBase err ", chainId, poolId, err)
			}
			_ = db.RedisSetString("base_info:lc_pool_"+chainId+"_"+poolId, baseInfoMd5Str, 60*30)
		}
		dataInfo, err := cli.PoolDataInfo(ctx, big.NewInt(int64(i)))
		if err != nil {
			log.Logger.Sugar().Info("UpdatePoolInfo PoolDataInfo err", poolId, err)
			continue
		}
		poolData := models.PoolData{
			PoolId:                 poolId,
			ChainId:                chainId,
			FinishAmountBorrow:     dataInfo.FinishAmountBorrow.String(),
			FinishAmountLend:       dataInfo.FinishAmountLend.String(),
			LiquidationAmounBorrow: dataInfo.LiquidationAmounBorrow.String(),
			LiquidationAmounLend:   dataInfo.LiquidationAmounLend.String(),
			SettleAmountBorrow:     dataInfo.SettleAmountBorrow.String(),
			SettleAmountLend:       dataInfo.SettleAmountLend.String(),
		}
		hasPoolData, byteDataInfoStr, dataInfoMd5Str := s.hashRedis("data_info:lc_pool_"+chainId+"_"+poolId, &poolData)
		if !hasPoolData || dataInfoMd5Str != byteDataInfoStr {
			err = models.NewPoolData().SavePoolData(chainId, poolId, &poolData)
			if err != nil {
				log.Logger.Sugar().Error("SavePoolData err ", chainId, poolId, err)
			}
			_ = db.RedisSetString("data_info:lc_pool_"+chainId+"_"+poolId, dataInfoMd5Str, 60*30)
		}
	}
}

func (s *poolService) GetPoolMd5(baseInfo *models.PoolBase, key string) (bool, string, string) {
	return s.hashRedis(key, baseInfo)
}

func (s *poolService) hashRedis(key string, v interface{}) (bool, string, string) {
	b, _ := json.Marshal(v)
	md5s := utils.Md5(string(b))
	resInfoBytes, _ := db.RedisGet(key)
	if len(resInfoBytes) > 0 {
		return true, strings.Trim(string(resInfoBytes), `"`), md5s
	}
	return false, "", md5s
}
