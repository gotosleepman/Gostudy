package services

import (
	"context"
	"math/big"

	"lending-copy/config"
	"lending-copy/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BalanceMonitor struct{}

func NewBalanceMonitor() *BalanceMonitor {
	return &BalanceMonitor{}
}

// Monitor 定时检查借贷合约地址原生币余额，低于阈值时打日志（pledge-backend 为邮件告警）
func (s *BalanceMonitor) Monitor() {
	net := config.Config.TestNet.NetUrl
	addr := common.HexToAddress(config.Config.TestNet.LendingPoolAddr)
	if addr == (common.Address{}) {
		return
	}
	bal, err := s.nativeBalance(net, addr)
	if err != nil {
		log.Logger.Sugar().Error("BalanceMonitor ", err)
		return
	}
	th, ok := new(big.Int).SetString(config.Config.Threshold.LendingPoolNativeThreshold, 10)
	if !ok {
		return
	}
	if bal.Cmp(th) <= 0 {
		log.Logger.Sugar().Warn("lending pool native balance low: contract=", addr.Hex(), " balance_wei=", bal.String(), " threshold=", th.String())
	}
}

func (s *BalanceMonitor) nativeBalance(netURL string, token common.Address) (*big.Int, error) {
	c, err := ethclient.Dial(netURL)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return c.BalanceAt(context.Background(), token, nil)
}
