package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"MAJOR-PROJECT/controllers"
	"MAJOR-PROJECT/routes"
	"MAJOR-PROJECT/util"

	"math/big"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Force IPv4 globally to fix S3/Network timeouts on IPv6-broken networks
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext(ctx, "tcp4", addr)
		}
	}

	// -----------------------------------------------------
	// 1) LOAD .ENV FIRST
	// -----------------------------------------------------
	if err := godotenv.Overload(); err != nil {
		log.Println("[WARN] No .env file found, using system environment variables")
	} else {
		log.Println("[OK] Loaded .env file successfully")
	}

	// -----------------------------------------------------
	// 2) DEPLOY CONTRACTS (IF NOT EXISTS)
	// -----------------------------------------------------
	privHex := os.Getenv("EVM_PRIVATE_KEY")
	if privHex == "" {
		log.Fatal("[ERROR] EVM_PRIVATE_KEY not set in .env. Execution aborted for security.")
	}

	// Deploy L2 Factory
	l2FactoryAddr := os.Getenv("L2_FACTORY_CONTRACT_ADDRESS")
	if l2FactoryAddr != "" {
		log.Printf("[OK] Existing L2 Factory found in .env: %s", l2FactoryAddr)
		log.Println("[SKIPPED] Skipping L2 Factory deployment.")
	} else {
		log.Println("[INFO] Deploying L2 Election Factory...")
		l2Rpc := os.Getenv("L2_NODE_URL")
		if l2Rpc == "" {
			log.Fatal("[ERROR] L2_NODE_URL must be set")
		}
		l2ChainIDStr := os.Getenv("L2_CHAIN_ID")
		l2ChainIDVal, _ := strconv.ParseInt(l2ChainIDStr, 10, 64)
		if l2ChainIDVal == 0 {
			l2ChainIDVal = 80002
		} // Default Amoy

		deployCtx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		addr, tx, err := util.Deploy(deployCtx, l2Rpc, "build/ElectionFact.abi", "build/ElectionFact.bin", privHex, big.NewInt(l2ChainIDVal))
		cancel()
		if err != nil {
			log.Fatalf("[ERROR] L2 Factory deployment failed: %v", err)
		}
		log.Printf("[L2 DEPLOYED] Factory contract at: %s | tx: %s", addr.Hex(), tx.Hash().Hex())
		os.Setenv("L2_FACTORY_CONTRACT_ADDRESS", addr.Hex())
	}

	// Deploy L1 Archive
	l1ArchiveAddr := os.Getenv("L1_ARCHIVE_CONTRACT_ADDRESS")
	if l1ArchiveAddr != "" {
		log.Printf("[OK] Existing L1 Archive found in .env: %s", l1ArchiveAddr)
		log.Println("[SKIPPED] Skipping L1 Archive deployment.")
	} else {
		log.Println("[INFO] Deploying L1 Election Archive...")
		l1Rpc := os.Getenv("L1_NODE_URL")
		if l1Rpc == "" {
			log.Fatal("[ERROR] L1_NODE_URL must be set")
		}
		l1ChainIDStr := os.Getenv("L1_CHAIN_ID")
		l1ChainIDVal, _ := strconv.ParseInt(l1ChainIDStr, 10, 64)
		if l1ChainIDVal == 0 {
			l1ChainIDVal = 11155111
		} // Default Sepolia

		deployCtx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		addr, tx, err := util.Deploy(deployCtx, l1Rpc, "build/Ethereum_Contract_ElectionArchive_sol_ElectionArchive.abi", "build/Ethereum_Contract_ElectionArchive_sol_ElectionArchive.bin", privHex, big.NewInt(l1ChainIDVal))
		cancel()
		if err != nil {
			log.Fatalf("[ERROR] L1 Archive deployment failed: %v", err)
		}
		log.Printf("[L1 DEPLOYED] Archive contract at: %s | tx: %s", addr.Hex(), tx.Hash().Hex())
		os.Setenv("L1_ARCHIVE_CONTRACT_ADDRESS", addr.Hex())
	}

	// Validate required env variables
	requiredEnvVars := []string{"MONGODB_URI", "EMAIL", "PASSWORD"}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			log.Printf("[WARN] Warning: Required environment variable %s is not set", v)
		}
	}

	// -----------------------------------------------------
	// 3) MONGO SETUP
	// -----------------------------------------------------
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	connectCtx, connectCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer connectCancel()

	clientOpts, tlsErr := buildMongoClientOptions(mongoURI)
	if tlsErr != nil {
		log.Fatalf("[ERROR] MongoDB TLS config error: %v", tlsErr)
	}
	client, err := mongo.Connect(connectCtx, clientOpts)
	if err != nil {
		log.Fatalf("[ERROR] MongoDB connect error: %v", err)
	}

	if err := client.Ping(connectCtx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		log.Fatalf("[ERROR] MongoDB ping error: %v", err)
	}

	fmt.Println("[OK] Connected to MongoDB successfully")

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "voting_system"
	}

	controllers.InitCompanyCollection(client, dbName)
	controllers.InitVoterCollection(client, dbName)
	controllers.InitCandidateCollection(client, dbName)
	controllers.InitOTPCollection(client, dbName)
	controllers.InitAuditCollection(client, dbName)
	controllers.InitMetadataCollection(client, dbName)
	controllers.InitStudentCollection(client, dbName)
	fmt.Println("[OK] Initialized database collections")

	// -----------------------------------------------------
	// 4) START SERVER
	// -----------------------------------------------------
	router := routes.SetupRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("[RUNNING] Server starting on port %s", port)
		log.Printf("[API] http://localhost:%s/api", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] Server error: %v", err)
		}
	}()

	// Wait for interrupt (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Shutting down server...")

	// Shutdown server gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[ERROR] Server forced to shutdown: %v", err)
	} else {
		log.Println("[OK] HTTP server stopped")
	}

	// Disconnect from MongoDB
	disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer disconnectCancel()

	if err := client.Disconnect(disconnectCtx); err != nil {
		log.Fatalf("[ERROR] MongoDB disconnect error: %v", err)
	}
	log.Println("[OK] MongoDB disconnected, exiting")
}

// buildMongoClientOptions adds TLS config for AWS DocumentDB when DOCDB_TLS_CA_FILE is set.
// Falls back to plain URI (Atlas / local) when the env var is absent.
// Also forces IPv4 dialing to avoid IPv6 NAT64 timeouts on some networks.
func buildMongoClientOptions(uri string) (*options.ClientOptions, error) {
	opts := options.Client().ApplyURI(uri)

	// Connection pooling for multi-day scaling limits
	opts.SetMaxPoolSize(100) // DocumentDB t3.medium can safely handle ~500. Cap the app at 100 to prevent exhaustion.
	opts.SetMinPoolSize(5)

	// Force IPv4 for mongo driver (its own dialer is separate from http.DefaultTransport)
	dialer := &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
	opts.SetDialer(&mongoDialer{dialer: dialer})

	caFile := os.Getenv("DOCDB_TLS_CA_FILE")
	if caFile == "" {
		return opts, nil
	}
	caPEM, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA file %s: %w", caFile, err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("failed to parse certificates from %s", caFile)
	}
	opts.SetTLSConfig(&tls.Config{
		RootCAs:            pool,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
	})
	log.Printf("[OK] DocumentDB TLS configured: %s", caFile)
	return opts, nil
}

// mongoDialer wraps net.Dialer to force IPv4 TCP connections for the mongo driver.
// Implements options.ContextDialer: DialContext(ctx, network, addr).
type mongoDialer struct{ dialer *net.Dialer }

func (d *mongoDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return d.dialer.DialContext(ctx, "tcp4", addr)
}
