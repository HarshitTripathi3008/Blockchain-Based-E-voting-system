# Blockchain-Based E-Voting System

A secure, transparent, and decentralized electronic voting system built with Go, MongoDB, and Ethereum blockchain technology. This system ensures the integrity of elections by recording key audit data on an immutable ledger.

## 🚀 Features

*   **Secure Voting:** leverages Ethereum blockchain to ensure vote records are immutable.
*   **Role-Based Access:** Distinct flows for Companies (Admins), Candidates, and Voters.
*   **Authentication:** Secure login and OTP verification.
*   **Real-time Dashboard:** Monitoring of election status and results.
*   **Audit Trail:** Comprehensive audit logs stored in MongoDB and hash references on Blockchain.
*   **Email Notifications:** Integrated email service for OTPs and election updates.

## 🛠️ Tech Stack

*   **Backend:** [Go (Golang)](https://golang.org/)
*   **Blockchain:** [Ethereum](https://ethereum.org/) (Smart Contracts)
*   **Database:** [MongoDB](https://www.mongodb.com/)
*   **Routing:** [Gorilla Mux](https://github.com/gorilla/mux)
*   **Ethereum Client:** [Go-Ethereum (geth)](https://geth.ethereum.org/)
*   **Frontend:** HTML/CSS/JS (served from `pages/` and `static/`)

## 📋 Prerequisites

Before running the application, ensure you have the following installed:

*   [Go](https://go.dev/dl/) (v1.24 or later)
*   [MongoDB](https://www.mongodb.com/try/download/community) (running locally or a cloud instance)
*   [Ganache](https://trufflesuite.com/ganache/) (or any Ethereum JSON-RPC node)
*   [Git](https://git-scm.com/)

## ⚙️ Installation & Setup

1.  **Clone the repository**
    ```bash
    git clone <repository-url>
    cd Blockchain-Based-E-voting-system
    ```

2.  **Install Go dependencies**
    ```bash
    go mod tidy
    ```

3.  **Configure Environment Variables**
    Create a `.env` file in the root directory. You can copy the following template:

    ```env
    # Server Configuration
    PORT=3000
    
    # Database Configuration
    MONGODB_URI=mongodb://localhost:27017
    DB_NAME=voting_system
    
    # Ethereum Configuration
    ETHEREUM_NODE_URL=http://127.0.0.1:8545
    ETHEREUM_PRIVATE_KEY=<YOUR_PRIVATE_KEY_FROM_GANACHE>
    ETHEREUM_CHAIN_ID=1337
    
    # Email Configuration
    EMAIL=<YOUR_EMAIL_FOR_SENDING_NOTIFICATIONS>
    PASSWORD=<YOUR_EMAIL_APP_PASSWORD>
    ```

4.  **Prepare Smart Contracts**
    Ensure the compiled strict ABI and BIN files for the contract are present in the `build/` directory:
    *   `build/ElectionFact.abi`
    *   `build/ElectionFact.bin`

    *(If these are missing, you may need to compile the Solidity contracts located in `Ethereum/Contract/` using `solc`.)*

## 🏃 Usage

1.  **Start MongoDB**
    Ensure your MongoDB service is running.

2.  **Start Ethereum Node**
    Open Ganache and ensure it is running on `127.0.0.1:8545`. Copy one of the private keys to your `.env` file.

3.  **Run the Application**
    ```bash
    go run main.go
    ```

    On startup, the application will:
    *   Connect to MongoDB.
    *   Deploy the `ElectionFact` smart contract to the local Ethereum network.
    *   Start the HTTP server.

4.  **Access the App**
    Open your browser and navigate to:
    [http://localhost:3000](http://localhost:3000)

## 📂 Project Structure

*   `controllers/`: Application logic and database interactions.
*   `routes/`: API route definitions.
*   `Ethereum/`: Smart contract sources.
*   `pages/`: HTML frontend pages.
*   `static/`: CSS, JS, and image assets.
*   `scripts/`: Utility scripts (e.g., clearing database data).
*   `util/`: Helper functions including Ethereum interactions.
