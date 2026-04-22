package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	rpc := os.Getenv("L2_NODE_URL")
	priv := os.Getenv("EVM_PRIVATE_KEY")

	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimSpace(priv))
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	id, _ := client.ChainID(context.Background())
	fmt.Printf("Chain ID: %v\n", id)

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatalf("Failed to retrieve balance: %v", err)
	}

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(1000000000000000000))

	fmt.Printf("Admin Address: %s\n", fromAddress.Hex())
	fmt.Printf("Balance: %s\n", ethValue.Text('f', 6))

	nonce, _ := client.PendingNonceAt(context.Background(), fromAddress)
	fmt.Printf("Pending Nonce: %d\n", nonce)
}
