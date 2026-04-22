package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	rpc := os.Getenv("L2_NODE_URL")
	if len(os.Args) < 2 {
		log.Fatal("Please provide contract address")
	}
	addr := os.Args[1]

	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}

	code, err := client.CodeAt(context.Background(), common.HexToAddress(addr), nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(code) > 0 {
		fmt.Printf("Contract %s exists on this network (%d bytes)\n", addr, len(code))
	} else {
		fmt.Printf("Contract %s DOES NOT EXIST on this network (EOA or empty)\n", addr)
	}
}
