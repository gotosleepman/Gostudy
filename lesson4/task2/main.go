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

	"lesson4/task2/bindings/counter"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const defaultChainID = 11155111

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "deploy":
		if err := runDeploy(os.Args[2:]); err != nil {
			log.Fatalf("部署合约失败: %v", err)
		}
	case "call":
		if err := runCall(os.Args[2:]); err != nil {
			log.Fatalf("调用合约失败: %v", err)
		}
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Counter 合约工具")
	fmt.Println("")
	fmt.Println("用法:")
	fmt.Println("  go run . deploy [-rpc <RPC_URL>] [-pk <PRIVATE_KEY>] [-chainid 11155111]")
	fmt.Println("  go run . call [-rpc <RPC_URL>] [-pk <PRIVATE_KEY>] [-addr <CONTRACT_ADDRESS>] [-chainid 11155111]")
	fmt.Println("")
	fmt.Println("环境变量:")
	fmt.Println("  RPC_URL")
	fmt.Println("  PRIVATE_KEY")
	fmt.Println("  CONTRACT_ADDRESS (仅 call 时可用)")
}

func runDeploy(args []string) error {
	fs := flag.NewFlagSet("deploy", flag.ContinueOnError)
	rpcURL := fs.String("rpc", "", "Sepolia RPC URL")
	privateKeyHex := fs.String("pk", "", "部署者私钥（16进制，支持带/不带 0x）")
	chainID := fs.Int64("chainid", defaultChainID, "链 ID，Sepolia 默认为 11155111")
	if err := fs.Parse(args); err != nil {
		return err
	}

	resolvedRPC := firstNonEmpty(*rpcURL, os.Getenv("RPC_URL"))
	resolvedPK := firstNonEmpty(*privateKeyHex, os.Getenv("PRIVATE_KEY"))
	if resolvedRPC == "" || resolvedPK == "" {
		return fmt.Errorf("缺少必要参数：需提供 RPC_URL、PRIVATE_KEY（命令行参数或环境变量）")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := ethclient.DialContext(ctx, resolvedRPC)
	if err != nil {
		return fmt.Errorf("连接 RPC 失败: %w", err)
	}
	defer client.Close()

	privKey, fromAddr, err := parsePrivateKey(resolvedPK)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(*chainID))
	if err != nil {
		return fmt.Errorf("创建签名器失败: %w", err)
	}
	auth.Context = ctx

	if err = fillDynamicFee(ctx, client, auth); err != nil {
		return fmt.Errorf("获取动态费参数失败: %w", err)
	}

	contractAddr, tx, _, err := counter.DeployCounter(auth, client)
	if err != nil {
		return fmt.Errorf("发送部署交易失败: %w", err)
	}
	fmt.Printf("部署交易已发送: %s\n", tx.Hash().Hex())

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("等待部署交易上链失败: %w", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("部署交易执行失败，txHash=%s", tx.Hash().Hex())
	}

	fmt.Println("=== 部署成功 ===")
	fmt.Printf("部署者: %s\n", fromAddr.Hex())
	fmt.Printf("合约地址: %s\n", contractAddr.Hex())
	fmt.Printf("区块高度: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("浏览器: https://sepolia.etherscan.io/tx/%s\n", tx.Hash().Hex())
	return nil
}

func runCall(args []string) error {
	fs := flag.NewFlagSet("call", flag.ContinueOnError)
	rpcURL := fs.String("rpc", "", "Sepolia RPC URL")
	privateKeyHex := fs.String("pk", "", "调用者私钥（16进制，支持带/不带 0x）")
	contractAddrHex := fs.String("addr", "", "已部署 Counter 合约地址")
	chainID := fs.Int64("chainid", defaultChainID, "链 ID，Sepolia 默认为 11155111")
	if err := fs.Parse(args); err != nil {
		return err
	}

	resolvedRPC := firstNonEmpty(*rpcURL, os.Getenv("RPC_URL"))
	resolvedPK := firstNonEmpty(*privateKeyHex, os.Getenv("PRIVATE_KEY"))
	resolvedAddr := firstNonEmpty(*contractAddrHex, os.Getenv("CONTRACT_ADDRESS"))

	if resolvedRPC == "" || resolvedPK == "" || resolvedAddr == "" {
		return fmt.Errorf("缺少必要参数：需提供 RPC_URL、PRIVATE_KEY、CONTRACT_ADDRESS（命令行参数或环境变量）")
	}
	if !common.IsHexAddress(resolvedAddr) {
		return fmt.Errorf("合约地址不合法: %s", resolvedAddr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := ethclient.DialContext(ctx, resolvedRPC)
	if err != nil {
		return fmt.Errorf("连接 RPC 失败: %w", err)
	}
	defer client.Close()

	contractAddr := common.HexToAddress(resolvedAddr)
	ctr, err := counter.NewCounter(contractAddr, client)
	if err != nil {
		return fmt.Errorf("创建合约实例失败: %w", err)
	}

	before, err := ctr.GetCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("读取调用前计数失败: %w", err)
	}
	fmt.Printf("调用前计数: %s\n", before.String())

	privKey, fromAddr, err := parsePrivateKey(resolvedPK)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(*chainID))
	if err != nil {
		return fmt.Errorf("创建签名器失败: %w", err)
	}
	auth.Context = ctx

	// 让绑定方法自动估算 gas；为了兼容 EIP-1559 手动填充费用参数。
	if err = fillDynamicFee(ctx, client, auth); err != nil {
		return fmt.Errorf("获取动态费参数失败: %w", err)
	}

	tx, err := ctr.Increment(auth)
	if err != nil {
		return fmt.Errorf("发送 increment 交易失败: %w", err)
	}
	fmt.Printf("交易已发送: %s\n", tx.Hash().Hex())

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("等待交易上链失败: %w", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("交易执行失败，txHash=%s", tx.Hash().Hex())
	}

	after, err := ctr.GetCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("读取调用后计数失败: %w", err)
	}

	fmt.Println("=== 调用结果 ===")
	fmt.Printf("发送方: %s\n", fromAddr.Hex())
	fmt.Printf("合约地址: %s\n", contractAddr.Hex())
	fmt.Printf("区块高度: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("调用后计数: %s\n", after.String())
	fmt.Printf("浏览器: https://sepolia.etherscan.io/tx/%s\n", tx.Hash().Hex())
	return nil
}

func parsePrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, common.Address, error) {
	privateKeyHex = strings.TrimPrefix(strings.TrimSpace(privateKeyHex), "0x")
	privKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("私钥格式错误: %w", err)
	}
	pubKey, ok := privKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, fmt.Errorf("解析公钥失败")
	}
	return privKey, crypto.PubkeyToAddress(*pubKey), nil
}

func fillDynamicFee(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts) error {
	tipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return err
	}
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}

	feeCap := new(big.Int).Mul(header.BaseFee, big.NewInt(2))
	feeCap.Add(feeCap, tipCap)

	auth.GasTipCap = tipCap
	auth.GasFeeCap = feeCap
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
