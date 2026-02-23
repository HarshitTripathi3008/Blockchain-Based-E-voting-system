package util

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Deploy deploys the contract and returns the deployed address and tx hash.
// It does NOT os.Exit / log.Fatalf; instead it returns errors to the caller.
func Deploy(ctx context.Context, rpc string, abiPath string, binPath string, privHex string, chainID *big.Int) (common.Address, *types.Transaction, error) {
	// Basic validation
	if rpc == "" {
		return common.Address{}, nil, errors.New("rpc URL is empty")
	}
	if abiPath == "" || binPath == "" {
		return common.Address{}, nil, errors.New("abiPath or binPath is empty")
	}
	if privHex == "" {
		return common.Address{}, nil, errors.New("private key is empty")
	}
	// dial with context and a short timeout derived from ctx
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(dialCtx, rpc)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to connect to RPC %s: %w", rpc, err)
	}
	// caller should close with client.Close() --- but ethclient doesn't export Close; its underlying connection will be cleaned by GC.
	// We'll continue to use the client as needed.

	// parse private key
	priv, err := crypto.HexToECDSA(strings.TrimPrefix(privHex, "0x"))
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("invalid private key: %w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(priv, chainID)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// Read and parse ABI
	abiBytes, err := os.ReadFile(abiPath)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to read ABI file %s: %w", abiPath, err)
	}
	parsedABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to parse ABI JSON: %w", err)
	}

	// Read and clean BIN
	binBytesRaw, err := os.ReadFile(binPath)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to read BIN file %s: %w", binPath, err)
	}
	binStr := strings.TrimSpace(string(binBytesRaw))
	binStr = strings.ReplaceAll(binStr, "\n", "")
	binStr = strings.ReplaceAll(binStr, "\r", "")
	binStr = strings.ReplaceAll(binStr, " ", "")
	if !strings.HasPrefix(binStr, "0x") {
		binStr = "0x" + binStr
	}
	bin := common.FromHex(binStr)

	// Deploy contract
	// Use provided ctx for waiting as well
	fmt.Println("[START] Deploying contract...")
	address, tx, _, err := bind.DeployContract(auth, parsedABI, bin, client)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to deploy contract: %w", err)
	}

	fmt.Printf("[TX] Deployment tx: %s\n", tx.Hash().Hex())
	fmt.Println("[WAIT] Waiting for deployment to be mined...")

	waitCtx, waitCancel := context.WithTimeout(ctx, 60*time.Second) // arbitrary wait timeout
	defer waitCancel()

	// WaitDeployed will poll the node until the contract is present or context times out
	_, err = bind.WaitDeployed(waitCtx, client, tx)
	if err != nil {
		return common.Address{}, tx, fmt.Errorf("failed waiting for deployment: %w", err)
	}

	fmt.Println("[OK] Deployed at address:", address.Hex())

	// Update or append FACTORY_CONTRACT_ADDRESS in .env (best-effort; if it fails, return the error)
	envPath := ".env"
	envData := ""
	if _, err := os.Stat(envPath); err == nil {
		// read existing .env
		b, err := os.ReadFile(envPath)
		if err != nil {
			return address, tx, fmt.Errorf("failed to read .env: %w", err)
		}
		envData = string(b)
	}
	lines := strings.Split(envData, "\n")
	found := false
	for i, l := range lines {
		if strings.HasPrefix(l, "FACTORY_CONTRACT_ADDRESS=") {
			lines[i] = "FACTORY_CONTRACT_ADDRESS=" + address.Hex()
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, "FACTORY_CONTRACT_ADDRESS="+address.Hex())
	}
	if err := os.WriteFile(envPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return address, tx, fmt.Errorf("failed to write .env: %w", err)
	}
	fmt.Println("[INFO] .env updated with FACTORY_CONTRACT_ADDRESS =", address.Hex())

	// Verify runtime code present at address
	codeCtx, codeCancel := context.WithTimeout(ctx, 10*time.Second)
	defer codeCancel()
	code, err := client.CodeAt(codeCtx, address, nil)
	if err != nil {
		return address, tx, fmt.Errorf("failed to fetch code at deployed address: %w", err)
	}
	if len(code) == 0 {
		fmt.Println("[WARN] Warning: runtime code at address is empty. Contract creation may have failed (invalid opcode, revert, etc.).")
	} else {
		fmt.Println("[OK] Runtime code stored at deployed address (bytes):", len(code))
	}

	return address, tx, nil
}
