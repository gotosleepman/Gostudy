package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	defaultChainID = 11155111
)

func usage() {
	fmt.Println("Sepolia 区块链交互工具")
	fmt.Println("")
	fmt.Println("用法:")
	fmt.Println("  go run . query-block [-rpc <RPC_URL>] -block <区块号>")
	fmt.Println("  go run . send-tx [-rpc <RPC_URL>] [-pk <私钥>] -to <接收地址> -amount <ETH数量>")
	fmt.Println("")
	fmt.Println("环境变量（推荐，避免敏感参数出现在命令历史）:")
	fmt.Println("  RPC_URL      Sepolia RPC URL")
	fmt.Println("  PRIVATE_KEY  发送方私钥（16进制，支持带/不带0x）")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  $env:RPC_URL='https://sepolia.infura.io/v3/<INFURA_KEY>'")
	fmt.Println("  $env:PRIVATE_KEY='<PRIVATE_KEY_HEX>'")
	fmt.Println("  go run . query-block -block 6000000")
	fmt.Println("  go run . send-tx -to 0xabc... -amount 0.0001")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "query-block":
		if err := runQueryBlock(os.Args[2:]); err != nil {
			log.Fatalf("查询区块失败: %v", err)
		}
	case "send-tx":
		if err := runSendTx(os.Args[2:]); err != nil {
			log.Fatalf("发送交易失败: %v", err)
		}
	default:
		usage()
		os.Exit(1)
	}
}

func runQueryBlock(args []string) error {
	fs := flag.NewFlagSet("query-block", flag.ContinueOnError)
	rpcURL := fs.String("rpc", "", "Sepolia RPC URL，例如 Infura endpoint")
	blockNum := fs.Uint64("block", 0, "区块号")
	if err := fs.Parse(args); err != nil {
		return err
	}
	resolvedRPC := firstNonEmpty(*rpcURL, os.Getenv("RPC_URL"))
	if resolvedRPC == "" {
		return fmt.Errorf("-rpc 为空且未设置环境变量 RPC_URL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, resolvedRPC)
	if err != nil {
		return fmt.Errorf("连接 RPC 失败: %w", err)
	}
	defer client.Close()

	block, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(*blockNum))
	if err != nil {
		return fmt.Errorf("获取区块失败: %w", err)
	}

	fmt.Println("=== 区块信息 ===")
	fmt.Printf("区块号: %d\n", block.NumberU64())
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("父区块哈希: %s\n", block.ParentHash().Hex())
	fmt.Printf("时间戳: %d (%s)\n", block.Time(), time.Unix(int64(block.Time()), 0).UTC().Format(time.RFC3339))
	fmt.Printf("交易数量: %d\n", len(block.Transactions()))
	fmt.Printf("矿工地址(coinbase): %s\n", block.Coinbase().Hex())
	fmt.Printf("GasUsed/GasLimit: %d / %d\n", block.GasUsed(), block.GasLimit())

	return nil
}

func runSendTx(args []string) error {
	fs := flag.NewFlagSet("send-tx", flag.ContinueOnError)
	rpcURL := fs.String("rpc", "", "Sepolia RPC URL，例如 Infura endpoint")
	privateKeyHex := fs.String("pk", "", "发送方私钥（16进制，不带0x）")
	toAddrHex := fs.String("to", "", "接收方地址")
	amountEth := fs.String("amount", "0", "转账金额（单位 ETH，例如 0.0001）")
	chainID := fs.Int64("chainid", defaultChainID, "链 ID，Sepolia 默认 11155111")
	if err := fs.Parse(args); err != nil {
		return err
	}
	resolvedRPC := firstNonEmpty(*rpcURL, os.Getenv("RPC_URL"))
	resolvedPrivateKey := firstNonEmpty(*privateKeyHex, os.Getenv("PRIVATE_KEY"))
	if resolvedRPC == "" || resolvedPrivateKey == "" || *toAddrHex == "" {
		return fmt.Errorf("缺少必要参数：需要 -to，且需提供 -rpc 或 RPC_URL、-pk 或 PRIVATE_KEY")
	}
	if !common.IsHexAddress(*toAddrHex) {
		return fmt.Errorf("接收地址格式不正确: %s", *toAddrHex)
	}

	privateKey, fromAddr, err := parsePrivateKey(resolvedPrivateKey)
	if err != nil {
		return err
	}
	toAddr := common.HexToAddress(*toAddrHex)

	valueWei, err := ethToWei(*amountEth)
	if err != nil {
		return fmt.Errorf("解析 -amount 失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, resolvedRPC)
	if err != nil {
		return fmt.Errorf("连接 RPC 失败: %w", err)
	}
	defer client.Close()

	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return fmt.Errorf("获取 nonce 失败: %w", err)
	}

	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return fmt.Errorf("获取 gasTipCap 失败: %w", err)
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("获取最新区块头失败: %w", err)
	}

	// EIP-1559: maxFee 设为 2*baseFee + tip，适合作为测试网络的保守估计
	gasFeeCap := new(big.Int).Mul(header.BaseFee, big.NewInt(2))
	gasFeeCap = gasFeeCap.Add(gasFeeCap, gasTipCap)

	const gasLimit = uint64(21000)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(*chainID),
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &toAddr,
		Value:     valueWei,
		Data:      nil,
	})

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(big.NewInt(*chainID)), privateKey)
	if err != nil {
		return fmt.Errorf("签名交易失败: %w", err)
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("广播交易失败: %w", err)
	}

	fmt.Println("=== 交易已发送 ===")
	fmt.Printf("发送方: %s\n", fromAddr.Hex())
	fmt.Printf("接收方: %s\n", toAddr.Hex())
	fmt.Printf("金额(wei): %s\n", valueWei.String())
	fmt.Printf("nonce: %d\n", nonce)
	fmt.Printf("交易哈希: %s\n", signedTx.Hash().Hex())
	fmt.Printf("浏览器链接: https://sepolia.etherscan.io/tx/%s\n", signedTx.Hash().Hex())

	return nil
}

func parsePrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, common.Address, error) {
	privateKeyHex = strings.TrimPrefix(strings.TrimSpace(privateKeyHex), "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("私钥格式不正确: %w", err)
	}
	pubKey := privateKey.Public()
	pubKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, fmt.Errorf("无法解析公钥")
	}
	fromAddr := crypto.PubkeyToAddress(*pubKeyECDSA)
	return privateKey, fromAddr, nil
}

func ethToWei(eth string) (*big.Int, error) {
	if strings.TrimSpace(eth) == "" {
		return nil, fmt.Errorf("金额为空")
	}
	r, ok := new(big.Rat).SetString(eth)
	if !ok {
		return nil, fmt.Errorf("无效金额: %s", eth)
	}

	weiBase := new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	r = new(big.Rat).Mul(r, weiBase)

	if !r.IsInt() {
		// 允许最多 18 位小数，超出则报错，避免精度隐患
		return nil, fmt.Errorf("金额精度超过 18 位小数: %s", eth)
	}
	return r.Num(), nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
