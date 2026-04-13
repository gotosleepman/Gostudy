package bindings

import (
	"context"
	_ "embed"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

//go:embed simple_lending.json
var abiJSON string

var parsed abi.ABI

func init() {
	var err error
	parsed, err = abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		panic(err)
	}
}

// PoolBaseTuple 与合约 poolBaseInfo 返回值一致
type PoolBaseTuple struct {
	SettleTime             *big.Int
	EndTime                *big.Int
	InterestRate           *big.Int
	MaxSupply              *big.Int
	LendSupply             *big.Int
	BorrowSupply           *big.Int
	MartgageRate           *big.Int
	LendToken              common.Address
	BorrowToken            common.Address
	State                  uint8
	SpCoin                 common.Address
	JpCoin                 common.Address
	AutoLiquidateThreshold *big.Int
}

// PoolDataTuple 与合约 poolDataInfo 返回值一致
type PoolDataTuple struct {
	SettleAmountLend       *big.Int
	SettleAmountBorrow     *big.Int
	FinishAmountLend       *big.Int
	FinishAmountBorrow     *big.Int
	LiquidationAmounLend   *big.Int
	LiquidationAmounBorrow *big.Int
}

// Client 轻量 RPC 封装（仿 pledge-backend 的 bindings 用法）
type Client struct {
	Eth        *ethclient.Client
	Contract   common.Address
	Close func()
}

func Dial(networkURL string, contractHex string) (*Client, error) {
	c, err := ethclient.Dial(networkURL)
	if err != nil {
		return nil, err
	}
	if !common.IsHexAddress(contractHex) {
		c.Close()
		return nil, fmt.Errorf("invalid contract address")
	}
	return &Client{Eth: c, Contract: common.HexToAddress(contractHex), Close: func() { c.Close() }}, nil
}

func (c *Client) call(ctx context.Context, data []byte) ([]byte, error) {
	msg := ethereum.CallMsg{To: &c.Contract, Data: data}
	return c.Eth.CallContract(ctx, msg, nil)
}

func (c *Client) LendFee(ctx context.Context) (*big.Int, error) {
	data, err := parsed.Pack("lendFee")
	if err != nil {
		return nil, err
	}
	out, err := c.call(ctx, data)
	if err != nil {
		return nil, err
	}
	vals, err := parsed.Unpack("lendFee", out)
	if err != nil {
		return nil, err
	}
	return vals[0].(*big.Int), nil
}

func (c *Client) BorrowFee(ctx context.Context) (*big.Int, error) {
	data, err := parsed.Pack("borrowFee")
	if err != nil {
		return nil, err
	}
	out, err := c.call(ctx, data)
	if err != nil {
		return nil, err
	}
	vals, err := parsed.Unpack("borrowFee", out)
	if err != nil {
		return nil, err
	}
	return vals[0].(*big.Int), nil
}

func (c *Client) PoolLength(ctx context.Context) (*big.Int, error) {
	data, err := parsed.Pack("poolLength")
	if err != nil {
		return nil, err
	}
	out, err := c.call(ctx, data)
	if err != nil {
		return nil, err
	}
	vals, err := parsed.Unpack("poolLength", out)
	if err != nil {
		return nil, err
	}
	return vals[0].(*big.Int), nil
}

func (c *Client) PoolBaseInfo(ctx context.Context, index *big.Int) (*PoolBaseTuple, error) {
	data, err := parsed.Pack("poolBaseInfo", index)
	if err != nil {
		return nil, err
	}
	out, err := c.call(ctx, data)
	if err != nil {
		return nil, err
	}
	vals, err := parsed.Unpack("poolBaseInfo", out)
	if err != nil {
		return nil, err
	}
	return &PoolBaseTuple{
		SettleTime:             vals[0].(*big.Int),
		EndTime:                vals[1].(*big.Int),
		InterestRate:           vals[2].(*big.Int),
		MaxSupply:              vals[3].(*big.Int),
		LendSupply:             vals[4].(*big.Int),
		BorrowSupply:           vals[5].(*big.Int),
		MartgageRate:           vals[6].(*big.Int),
		LendToken:              vals[7].(common.Address),
		BorrowToken:            vals[8].(common.Address),
		State:                  vals[9].(uint8),
		SpCoin:                 vals[10].(common.Address),
		JpCoin:                 vals[11].(common.Address),
		AutoLiquidateThreshold: vals[12].(*big.Int),
	}, nil
}

func (c *Client) PoolDataInfo(ctx context.Context, index *big.Int) (*PoolDataTuple, error) {
	data, err := parsed.Pack("poolDataInfo", index)
	if err != nil {
		return nil, err
	}
	out, err := c.call(ctx, data)
	if err != nil {
		return nil, err
	}
	vals, err := parsed.Unpack("poolDataInfo", out)
	if err != nil {
		return nil, err
	}
	return &PoolDataTuple{
		SettleAmountLend:       vals[0].(*big.Int),
		SettleAmountBorrow:     vals[1].(*big.Int),
		FinishAmountLend:       vals[2].(*big.Int),
		FinishAmountBorrow:     vals[3].(*big.Int),
		LiquidationAmounLend:   vals[4].(*big.Int),
		LiquidationAmounBorrow: vals[5].(*big.Int),
	}, nil
}
