# Blockchain-Based E-Voting System (SecureVote)

SecureVote is a next-generation electronic voting platform that combines the speed and cost-efficiency of Layer 2 (L2) with the ultimate security and immutability of Layer 1 (L1) Ethereum anchoring.

## 🏗️ Technical Architecture (2-Chain System)

The system operates on an advanced **Dual-Chain Distributed Ledger** architecture:

### 1. Active Voting Layer (Layer 2 - Polygon Amoy)
*   **Purpose:** High-throughput, low-latency execution of voting activities.
*   **Contract:** `ElectionFact.sol` (Factory Pattern).
*   **Functionality:** 
    *   Companies (Admins) can deploy their own dedicated `Election` smart contracts.
    *   Registers candidates and voters on-chain.
    *   Handles real-time vote casting with minimal gas fees.
    *   Calculates real-time winners and turnout statistics.

### 2. Result Anchoring Layer (Layer 1 - Ethereum Sepolia)
*   **Purpose:** Permanent, high-security archival of final results.
*   **Contract:** `ElectionArchive.sol`.
*   **Process:** 
    *   When an admin ends an election, the backend executes an **Anchoring Process**.
    *   It fetches final tallies and winner IDs from the L2 contract.
    *   It pushes a signed transaction to the L1 Archive contract, "locking" the results on the main Ethereum testnet.
    *   This ensures that even if L2 data were theoretically compromised, the final verified result remains immutable on L1.

---

## 🔒 Security & Authorization

To protect against unauthorized access and URL manipulation, SecureVote implements a multi-layered security model:

*   **Dashboard Gatekeeping:** All dashboards (`Admin`, `Voter`, `Observer`) utilize cookie-based JavaScript redirects. If the required session cookie (e.g., `company_email`, `voter_email`, or `observer_mode`) is missing or invalid, the user is immediately redirected to the portal login before any data is rendered.
*   **OTP Verification:** Voter authentication is hardened with 2-Factor Authentication (OTP) sent via secure email channels.
*   **Audit Logging:** Every critical action (Election Start, Vote Cast, Election End, L1 Anchoring) is logged in a centralized MongoDB Audit Trail and reference-hashed periodically.
*   **Image Integrity:** Candidate and Voter photos are stored on **AWS S3** with high availability and served over secure channels.

---

## ⚙️ Setup & Installation Instructions

If you are pulling this repository for the first time, follow these steps to get the system running locally or in production.

### 1. Prerequisites
*   **Go:** v1.21 or later.
*   **MongoDB:** A running instance (Local, MongoDB Atlas, or AWS DocumentDB).
*   **Wallet:** An Ethereum private key with testnet funds (Sepolia ETH and Amoy POL).

### 2. Installation
```bash
git clone https://github.com/HarshitTripathi3008/Blockchain-Based-E-voting-system.git
cd Blockchain-Based-E-voting-system
go mod tidy
```

### 3. Environment Configuration
Create a `.env` file in the root directory and populate it with your credentials. 

> [!IMPORTANT]
> You only need to provide the **RPC URLs** and **Private Key**. The system is designed to **auto-deploy** the necessary smart contracts and update the `.env` file with the contract addresses on the first run.

```env
PORT=3000
MONGODB_URI=mongodb://localhost:27017
DB_NAME=voting_system

# Shared Admin Wallet
EVM_PRIVATE_KEY=your_private_key_here

# Layer 1 (Sepolia) Configuration
L1_NODE_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_KEY
L1_CHAIN_ID=11155111

# Layer 2 (Amoy) Configuration
L2_NODE_URL=https://polygon-amoy.g.alchemy.com/v2/YOUR_KEY
L2_CHAIN_ID=80002

# Services
EMAIL=your-system-email@gmail.com
PASSWORD=your-app-specific-password
AWS_S3_BUCKET=your-secure-bucket-name
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
```

### 4. Running the Application
```bash
go run main.go
```
**What happens on startup?**
1.  Loads `.env`.
2.  **Auto-Deployment:** If `L2_FACTORY_CONTRACT_ADDRESS` is missing, it deploys the Factory to Polygon Amoy.
3.  **Auto-Deployment:** If `L1_ARCHIVE_CONTRACT_ADDRESS` is missing, it deploys the Archive to Ethereum Sepolia.
4.  Updates your `.env` with the new addresses.
5.  Connects to MongoDB and starts the HTTP server on port 3000.

---

## 🚀 Deployment (AWS Production)

The project includes optimized shell scripts for AWS EC2 deployment:
*   `setup_ec2.sh`: Installs Go, PM2, and Nginx.
*   `deploy_ec2.sh`: Automates `git pull`, build, and `pm2` restart.
*   `setup_ssl.sh`: Configures Certbot for HTTPS.

To update the live server:
```bash
ssh -i key.pem ubuntu@your-ip "cd ~/MAJOR_PROJECT && git pull && pm2 restart evoting"
```

---

## 📂 Project Structure

*   `controllers/`: Logic for L1 result anchoring, voter ID generation, and election management.
*   `bindings/`: Go-Ethereum generated bindings for smart contracts.
*   `Ethereum/Contract/`: Solidity source files for L1 and L2.
*   `pages/`: Secured frontend dashboards (Vanilla JS + Glassmorphism CSS).
*   `middleware/`: Server-side request filtering (IPv4 forcing, Logging).
*   `util/`: Deployment engine and blockchain transaction helpers.
