package controllers

import (
    "context"
    "encoding/json"
    "os"
    "net/http"
    "time"
    "MAJOR-PROJECT/bindings"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "golang.org/x/crypto/bcrypt"
)

type Company struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Email    string             `bson:"email" json:"email"`
    Password string             `bson:"password" json:"-"`
}

type CompanyRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type CompanyResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

var companyCollection *mongo.Collection

// Initialize company collection and unique email index
func InitCompanyCollection(client *mongo.Client, dbName string) {
    companyCollection = client.Database(dbName).Collection("companies")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _, _ = companyCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys:    bson.D{{Key: "email", Value: 1}},
        Options: options.Index().SetUnique(true),
    })
}

func withCompanyCORS(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func CreateCompany(w http.ResponseWriter, r *http.Request) {
    withCompanyCORS(w)
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req CompanyRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(CompanyResponse{Status: "error", Message: "Invalid request body"})
        return
    }
    if req.Email == "" || req.Password == "" {
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(CompanyResponse{Status: "error", Message: "email and password are required"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        _ = json.NewEncoder(w).Encode(CompanyResponse{Status: "error", Message: "Error hashing password"})
        return
    }

    newCompany := Company{Email: req.Email, Password: string(hashedPassword)}
    result, err := companyCollection.InsertOne(ctx, newCompany)
    if err != nil {
        if mongoErr, ok := err.(mongo.WriteException); ok {
            for _, we := range mongoErr.WriteErrors {
                if we.Code == 11000 {
                    w.WriteHeader(http.StatusConflict)
                    _ = json.NewEncoder(w).Encode(CompanyResponse{Status: "error", Message: "Company already exists"})
                    return
                }
            }
        }
        w.WriteHeader(http.StatusInternalServerError)
        _ = json.NewEncoder(w).Encode(CompanyResponse{Status: "error", Message: "Error creating company"})
        return
    }

    var created Company
    if err := companyCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&created); err != nil {
        created.ID = result.InsertedID.(primitive.ObjectID)
        created.Email = req.Email
    }

    _ = json.NewEncoder(w).Encode(CompanyResponse{
        Status:  "success",
        Message: "Company added successfully!!!",
        Data:    map[string]interface{}{"id": created.ID.Hex(), "email": created.Email},
    })
}

func AuthenticateCompany(w http.ResponseWriter, r *http.Request) {
    withCompanyCORS(w)
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req CompanyRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(CompanyResponse{Status: "error", Message: "Invalid request body"})
        return
    }
    if req.Email == "" || req.Password == "" {
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(CompanyResponse{Status: "error", Message: "email and password are required"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var companyInfo Company
    if err := companyCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&companyInfo); err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        _ = json.NewEncoder(w).Encode(CompanyResponse{Status: "error", Message: "Invalid email/password!!!"})
        return
    }

    if bcrypt.CompareHashAndPassword([]byte(companyInfo.Password), []byte(req.Password)) != nil {
        w.WriteHeader(http.StatusUnauthorized)
        _ = json.NewEncoder(w).Encode(CompanyResponse{Status: "error", Message: "Invalid email/password!!!"})
        return
    }

    // Try to find deployed election address from the factory contract (optional)
    electionAddrHex := ""
    factoryAddrStr := os.Getenv("FACTORY_CONTRACT_ADDRESS")
    if factoryAddrStr != "" {
        // getClient is defined in election_controller.go and is in the same package
        client, err := getClient()
        if err == nil {
            // ensure client closed
            defer client.Close()

            factoryAddr := common.HexToAddress(factoryAddrStr)
            factory, err := bindings.NewElectionFactory(factoryAddr, client)
            if err == nil {
                // Use call opts with request context (non-blocking/polite)
                callOpts := &bind.CallOpts{Context: r.Context(), Pending: false}
                deployedAddr, _, _, err := factory.GetDeployedElection(callOpts, req.Email)
                if err == nil {
                    if deployedAddr != (common.Address{}) {
                        electionAddrHex = deployedAddr.Hex()
                    }
                }
                // ignore errors — just return login without election address
            }
        }
    }

    data := map[string]interface{}{"id": companyInfo.ID.Hex(), "email": companyInfo.Email}
    if electionAddrHex != "" {
        data["election_address"] = electionAddrHex
    }

    _ = json.NewEncoder(w).Encode(CompanyResponse{
        Status:  "success",
        Message: "company found!!!",
        Data:    data,
    })
}
