package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"MAJOR-PROJECT/controllers"
	"MAJOR-PROJECT/routes"
	"MAJOR-PROJECT/util"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/big"
)

func main() {

	// -----------------------------------------------------
	// 1) DEPLOY CONTRACT BEFORE LOADING .ENV
	// -----------------------------------------------------
	{
		deployCtx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		rpc := "http://127.0.0.1:8545" 
		abiPath := "build/ElectionFact.abi"
		binPath := "build/ElectionFact.bin"
		privHex := "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d"
		chainID := big.NewInt(1337)

		addr, tx, err := util.Deploy(deployCtx, rpc, abiPath, binPath, privHex, chainID)
		if err != nil {
			log.Fatalf("❌ Contract deployment failed: %v", err)
		}
		log.Printf("🏗️ Contract deployed at: %s | tx: %s", addr.Hex(), tx.Hash().Hex())
	}

	// -----------------------------------------------------
	// 2) LOAD .ENV AFTER DEPLOY UPDATED IT
	// -----------------------------------------------------
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	} else {
		log.Println("✅ Loaded .env file successfully (after deployment)")
	}

	// Validate required env variables
	requiredEnvVars := []string{"MONGODB_URI", "EMAIL", "PASSWORD"}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			log.Printf("⚠️ Warning: Required environment variable %s is not set", v)
		}
	}

	// -----------------------------------------------------
	// 3) MONGO SETUP
	// -----------------------------------------------------
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	connectCtx, connectCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer connectCancel()

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(connectCtx, clientOpts)
	if err != nil {
		log.Fatalf("❌ MongoDB connect error: %v", err)
	}

	if err := client.Ping(connectCtx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		log.Fatalf("❌ MongoDB ping error: %v", err)
	}

	fmt.Println("✅ Connected to MongoDB successfully")

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "voting_system"
	}

	controllers.InitCompanyCollection(client, dbName)
	controllers.InitVoterCollection(client, dbName)
	controllers.InitCandidateCollection(client, dbName)
	controllers.InitOTPCollection(client, dbName)
	fmt.Println("✅ Initialized database collections")

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
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		log.Printf("📡 API: http://localhost:%s/api", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Wait for interrupt (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Shutdown server gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	} else {
		log.Println("✅ HTTP server stopped")
	}

	// Disconnect from MongoDB
	disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer disconnectCancel()

	if err := client.Disconnect(disconnectCtx); err != nil {
		log.Fatalf("❌ MongoDB disconnect error: %v", err)
	}
	log.Println("✅ MongoDB disconnected, exiting")
}
