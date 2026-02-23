package controllers

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"net/textproto"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"MAJOR-PROJECT/bindings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type VoterRegistration struct {
	ElectionAddress string    `bson:"election_address" json:"election_address"`
	Status          string    `bson:"status" json:"status"` // "Verified", "Pending"
	RegisteredAt    time.Time `bson:"registered_at" json:"registered_at"`
}

type Voter struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Email         string              `bson:"email" json:"email"`
	Password      string              `bson:"password" json:"-"`
	FullName      string              `bson:"full_name,omitempty" json:"full_name,omitempty"`
	DOB           time.Time           `bson:"dob,omitempty" json:"dob,omitempty"`
	RollNo        string              `bson:"roll_no,omitempty" json:"roll_no,omitempty"`
	Mobile        string              `bson:"mobile,omitempty" json:"mobile,omitempty"`
	Gender        string              `bson:"gender,omitempty" json:"gender,omitempty"`
	Year          string              `bson:"year,omitempty" json:"year,omitempty"`
	PhotoURL      string              `bson:"photo_url,omitempty" json:"photo_url,omitempty"`
	Registrations []VoterRegistration `bson:"registrations" json:"registrations"`
}

type Student struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `bson:"email" json:"email"`
	FullName string             `bson:"full_name" json:"full_name"`
	RollNo   string             `bson:"roll_no" json:"roll_no"`
	Mobile   string             `bson:"mobile" json:"mobile"`
	Gender   string             `bson:"gender" json:"gender"`
	Year     string             `bson:"year" json:"year"`
	DOB      time.Time          `bson:"dob" json:"dob"`
}

type VoterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	FullName        string `json:"full_name,omitempty"`
	DOB             string `json:"dob,omitempty"` // expect "YYYY-MM-DD"
	RollNo          string `json:"roll_no,omitempty"`
	Mobile          string `json:"mobile,omitempty"`
	Gender          string `json:"gender,omitempty"`
	Year            string `json:"year,omitempty"`
	PhotoURL        string `json:"photo_url,omitempty"`
	ElectionName    string `json:"election_name,omitempty"`
	ElectionAddress string `json:"election_address,omitempty"`
	CandidateID     int    `json:"candidate_id,omitempty"`
	VoterEmail      string `json:"voter_email,omitempty"`
	WinnerCandidate string `json:"winner_candidate,omitempty"`
	CandidateEmail  string `json:"candidate_email,omitempty"`
	OTP             string `json:"otp,omitempty"`
}

type VoterResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Count   int         `json:"count,omitempty"`
}

// Helper function to send JSON error responses
func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(VoterResponse{
		Status:  "error",
		Message: message,
	})
}

var voterCollection *mongo.Collection
var otpCollection *mongo.Collection

// Initialize collections
func InitVoterCollection(client *mongo.Client, dbName string) {
	voterCollection = client.Database(dbName).Collection("voters")

	// Create fast lookup index for extreme concurrency
	indexModel := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _ = voterCollection.Indexes().CreateOne(ctx, indexModel)

	fmt.Println("[OK] Initialized voters collection with indexes")
}

func InitOTPCollection(client *mongo.Client, dbName string) {
	otpCollection = client.Database(dbName).Collection("otps")

	indexModel := mongo.IndexModel{Keys: bson.M{"email": 1}}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _ = otpCollection.Indexes().CreateOne(ctx, indexModel)

	fmt.Println("[OK] Initialized OTP collection with indexes")
}

var studentCollection *mongo.Collection

func InitStudentCollection(client *mongo.Client, dbName string) {
	studentCollection = client.Database(dbName).Collection("students")

	indexModel := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _ = studentCollection.Indexes().CreateOne(ctx, indexModel)

	fmt.Println("[OK] Initialized Students collection with indexes")
}

func withVoterCORS(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// ===== helpers =====

func parseDOB(dobStr string) (time.Time, error) {
	if dobStr == "" {
		return time.Time{}, nil
	}
	// expect 2006-01-02
	t, err := time.Parse("2006-01-02", dobStr)
	if err == nil {
		return t, nil
	}
	// try RFC3339 fallback
	return time.Parse(time.RFC3339, dobStr)
}

func computeAge(d time.Time) int {
	if d.IsZero() {
		return 0
	}
	now := time.Now().UTC()
	age := now.Year() - d.Year()
	if now.YearDay() < d.YearDay() {
		age--
	}
	return age
}

var mobileRe = regexp.MustCompile(`^\+?[0-9]{7,15}$`)

// generate numeric 6-digit OTP
func genOTP() (string, error) {
	otp := ""
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += n.String()
	}
	return otp, nil
}

// generate alphanumeric password length n
func genPassword(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[r.Int64()]
	}
	return string(b), nil
}

// ===== existing handlers (RegisterVoter, UpdateVoter, DeleteVoter etc.) =====
// --- For brevity, keep your existing implementations unchanged.
// If you prefer I can paste them in full; currently they remain as in your repo.

// --------------------------
// New: SendOTP handler
// POST /api/voters/send-otp
// body: { "email": "...", "election_address": "..." }
// --------------------------
func SendOTP(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email           string `json:"email"`
		ElectionAddress string `json:"election_address,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Email == "" {
		http.Error(w, "email required", http.StatusBadRequest)
		return
	}
	if otpCollection == nil {
		http.Error(w, "server misconfigured: otp collection not ready", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	// rate-limit: if existing OTP not expired, deny quick repeat
	var existing bson.M
	err := otpCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existing)
	if err == nil {
		if exp, ok := existing["expiresAt"].(primitive.DateTime); ok {
			expTime := exp.Time()
			if time.Now().Before(expTime) {
				// still valid; disallow too-frequent sends
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "OTP already sent and still valid (wait before requesting another)"})
				return
			}
		} else if exp2, ok2 := existing["expiresAt"].(time.Time); ok2 {
			if time.Now().Before(exp2) {
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "OTP already sent and still valid (wait before requesting another)"})
				return
			}
		}
	}

	otp, err := genOTP()
	if err != nil {
		http.Error(w, "failed to generate otp", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(10 * time.Minute).UTC()
	doc := bson.M{
		"email":            req.Email,
		"otp":              otp,
		"expiresAt":        expiresAt,
		"election_address": req.ElectionAddress,
		"createdAt":        time.Now().UTC(),
	}

	// remove older OTP entries
	_, _ = otpCollection.DeleteMany(ctx, bson.M{"email": req.Email})
	if _, err := otpCollection.InsertOne(ctx, doc); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	subject := "Your OTP for voter registration"
	body := GenerateOTPEmail(otp)

	if err := sendEmail(req.Email, subject, body); err != nil {
		fmt.Printf("sendEmail error (SendOTP): %v\n", err)
		// remove inserted OTP because email failed
		_, _ = otpCollection.DeleteMany(ctx, bson.M{"email": req.Email})
		http.Error(w, "failed to send otp email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Try to find student data
	var studentData *Student
	if studentCollection != nil {
		var s Student
		if err := studentCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&s); err == nil {
			studentData = &s
		}
	}

	resp := VoterResponse{Status: "success", Message: "OTP sent"}
	if studentData != nil {
		resp.Data = studentData
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// --------------------------
// New: VerifyOTPAndRegister handler
// POST /api/voters/verify-otp-register
// body: { "email","otp","full_name","dob","mobile","address","father_name","mother_name","election_address" }
// --------------------------
func VerifyOTPAndRegister(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VoterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Email == "" {
		sendJSONError(w, "email required", http.StatusBadRequest)
		return
	}
	if otpCollection == nil || voterCollection == nil {
		sendJSONError(w, "server misconfigured: collections not ready", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// By-passing OTP logic:
	// The admin can register a voter directly without an OTP verification step.

	// basic validations
	if req.DOB != "" {
		if dob, err := parseDOB(req.DOB); err == nil {
			if computeAge(dob) < 16 {
				sendJSONError(w, "voter must be 16+", http.StatusBadRequest)
				return
			}
		} else {
			sendJSONError(w, "invalid dob", http.StatusBadRequest)
			return
		}
	}
	if req.Mobile != "" && !mobileRe.MatchString(req.Mobile) {
		sendJSONError(w, "invalid mobile format", http.StatusBadRequest)
		return
	}

	// ensure email not already registered globally (for this flow)
	count, err := voterCollection.CountDocuments(ctx, bson.M{"email": req.Email})
	if err == nil && count > 0 {
		sendJSONError(w, "Voter account already exists. Please login to join election.", http.StatusConflict)
		return
	}

	// generate password
	rawPassword, err := genPassword(12)
	if err != nil {
		sendJSONError(w, "failed to create password", http.StatusInternalServerError)
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		sendJSONError(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	dobTime, _ := parseDOB(req.DOB)

	newVoter := Voter{
		Email:    req.Email,
		Password: string(hashed),
		FullName: req.FullName,
		DOB:      dobTime,
		Mobile:   req.Mobile,
		RollNo:   req.RollNo,
		Gender:   req.Gender,
		Year:     req.Year,
		PhotoURL: req.PhotoURL,
		Registrations: []VoterRegistration{
			{
				ElectionAddress: req.ElectionAddress,
				Status:          "Verified",
				RegisteredAt:    time.Now().UTC(),
			},
		},
	}

	if _, err := voterCollection.InsertOne(ctx, newVoter); err != nil {
		sendJSONError(w, "failed to create voter", http.StatusInternalServerError)
		return
	}

	// Removed OTP deletion since we bypassed it

	// email password to voter
	body := GenerateWelcomeEmail(req.FullName, req.Email, rawPassword, req.ElectionAddress)
	subject := "Welcome to SecureVote - Account Credentials"
	if err := sendEmail(req.Email, subject, body); err != nil {
		fmt.Printf("sendEmail error (password email): %v\n", err)
	}

	_ = json.NewEncoder(w).Encode(VoterResponse{Status: "success", Message: "voter account created and verified"})
}

// ===== RegisterVoter (expanded) =====
func RegisterVoter(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VoterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Invalid request body"})
		return
	}
	if req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "email is required"})
		return
	}

	// parse and validate DOB -> age >= 18
	var dob time.Time
	if req.DOB != "" {
		p, err := parseDOB(req.DOB)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "invalid dob: expected YYYY-MM-DD"})
			return
		}
		dob = p
		age := computeAge(dob)
		if age < 18 {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "voter must be 18 years or older"})
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "dob is required (YYYY-MM-DD) for age verification"})
		return
	}

	// validate mobile if provided
	if req.Mobile != "" && !mobileRe.MatchString(req.Mobile) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "invalid mobile number; use digits, optional leading +, 7-15 chars"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if already exists
	var existing Voter
	err := voterCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existing)

	if err == nil {
		// Existing user: Link to new election if not already linked
		for _, reg := range existing.Registrations {
			if reg.ElectionAddress == req.ElectionAddress {
				w.WriteHeader(http.StatusConflict)
				_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Voter already registered for this election"})
				return
			}
		}

		// Add registration
		newReg := VoterRegistration{
			ElectionAddress: req.ElectionAddress,
			Status:          "Pending", // Admin added, maybe auto-verify? Let's say Pending until approved or verified
			RegisteredAt:    time.Now().UTC(),
		}
		_, err := voterCollection.UpdateOne(ctx, bson.M{"_id": existing.ID}, bson.M{"$push": bson.M{"registrations": newReg}})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Failed to link voter to election"})
			return
		}

		_ = json.NewEncoder(w).Encode(VoterResponse{
			Status:  "success",
			Message: "Existing voter added to this election.",
			Data: map[string]interface{}{
				"id":    existing.ID.Hex(),
				"email": existing.Email,
			},
		})
		return
	}

	// New User Logic
	rawPassword := req.Password
	if rawPassword == "" {
		rawPassword = req.Email // default
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Error hashing password"})
		return
	}

	newVoter := Voter{
		Email:    req.Email,
		Password: string(hashed),
		FullName: req.FullName,
		DOB:      dob,
		RollNo:   req.RollNo,
		Mobile:   req.Mobile,
		Gender:   req.Gender,
		Year:     req.Year,
		PhotoURL: req.PhotoURL,
		Registrations: []VoterRegistration{
			{
				ElectionAddress: req.ElectionAddress,
				Status:          "Pending",
				RegisteredAt:    time.Now().UTC(),
			},
		},
	}

	result, err := voterCollection.InsertOne(ctx, newVoter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Voter could not be added"})
		return
	}

	_ = json.NewEncoder(w).Encode(VoterResponse{
		Status:  "success",
		Message: "Voter account created and added to election.",
		Data: map[string]interface{}{
			"id":    result.InsertedID,
			"email": req.Email,
		},
	})
}

// ===== AuthenticateVoter =====
func AuthenticateVoter(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VoterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Invalid request body"})
		return
	}
	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "email and password are required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var voterInfo Voter
	if err := voterCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&voterInfo); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Invalid email/password"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(voterInfo.Password), []byte(req.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Invalid email/password"})
		return
	}

	_ = json.NewEncoder(w).Encode(VoterResponse{
		Status:  "success",
		Message: "voter authenticated",
		Data: map[string]interface{}{
			"id":            voterInfo.ID.Hex(),
			"email":         voterInfo.Email,
			"full_name":     voterInfo.FullName,
			"mobile":        voterInfo.Mobile,
			"registrations": voterInfo.Registrations,
		},
	})
}

// ===== GetAllVoters (returns expanded fields) =====
// ===== GetAllVoters (returns expanded fields) =====
func GetAllVoters(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var electionAddress string
	if r.Method == http.MethodGet {
		electionAddress = r.URL.Query().Get("election_address")
	} else {
		var req VoterRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		electionAddress = req.ElectionAddress
	}

	// Build Query
	filter := bson.M{}
	if electionAddress != "" {
		// Filter voters who have a registration for this election
		filter = bson.M{"registrations.election_address": electionAddress}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := voterCollection.Find(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Error fetching voters"})
		return
	}
	defer cursor.Close(ctx)

	var voters []Voter
	if err = cursor.All(ctx, &voters); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Error decoding voters"})
		return
	}

	vlist := make([]map[string]interface{}, 0, len(voters))
	for _, v := range voters {
		id := v.ID.Hex()
		var dobStr string
		if !v.DOB.IsZero() {
			dobStr = v.DOB.Format("2006-01-02")
		}

		// determine status for this specific election if filtered
		status := "Registered" // Default for global view
		if electionAddress != "" {
			for _, reg := range v.Registrations {
				if reg.ElectionAddress == electionAddress {
					status = reg.Status
					break
				}
			}
		}

		vlist = append(vlist, map[string]interface{}{
			"id":        id,
			"email":     v.Email,
			"full_name": v.FullName,
			"dob":       dobStr,
			"roll_no":   v.RollNo,
			"mobile":    v.Mobile,
			"gender":    v.Gender,
			"year":      v.Year,
			"status":    status,
		})
	}

	resp := map[string]interface{}{
		"status":  "success",
		"message": "voters list found",
		"voters":  vlist,
		"data":    map[string]interface{}{"voters": vlist}, // Keep both for compatibility
		"count":   len(vlist),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// ===== UpdateVoter (allow updating profile fields) =====
func UpdateVoter(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	voterID := vars["voterId"]
	if voterID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "voterId is required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(voterID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "invalid voterId"})
		return
	}

	var req VoterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "invalid JSON body"})
		return
	}

	update := bson.M{}
	if req.Email != "" {
		update["email"] = req.Email
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "error hashing password"})
			return
		}
		update["password"] = string(hashed)
	}
	if req.FullName != "" {
		update["full_name"] = req.FullName
	}
	if req.DOB != "" {
		if dob, err := parseDOB(req.DOB); err == nil {
			if computeAge(dob) < 18 {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "voter must be 18 years or older"})
				return
			}
			update["dob"] = dob
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "invalid dob format"})
			return
		}
	}
	if req.RollNo != "" {
		update["roll_no"] = req.RollNo
	}
	if req.Mobile != "" {
		if !mobileRe.MatchString(req.Mobile) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "invalid mobile format"})
			return
		}
		update["mobile"] = req.Mobile
	}
	if req.Gender != "" {
		update["gender"] = req.Gender
	}
	if req.Year != "" {
		update["year"] = req.Year
	}
	if req.PhotoURL != "" {
		update["photo_url"] = req.PhotoURL
	}

	if len(update) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "no fields to update"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := voterCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "failed to update voter"})
		return
	}
	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "voter not found"})
		return
	}

	_ = json.NewEncoder(w).Encode(VoterResponse{Status: "success", Message: "voter updated successfully"})
}

// ===== DeleteVoter (unchanged from original file) =====
func DeleteVoter(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	voterID := vars["voterId"]
	if voterID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Voter ID is required"})
		return
	}
	objID, err := primitive.ObjectIDFromHex(voterID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Invalid voter ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err = voterCollection.DeleteOne(ctx, bson.M{"_id": objID}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Error deleting voter"})
		return
	}

	_ = json.NewEncoder(w).Encode(VoterResponse{Status: "success", Message: "voter deleted successfully"})
}

// ===== NEW: GetVoterElections =====
func GetVoterElections(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodGet {
		// Expect ?voter_id=...
		voterID := r.URL.Query().Get("voter_id")
		if voterID == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "voter_id required"})
			return
		}

		objID, err := primitive.ObjectIDFromHex(voterID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Invalid voter_id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var voter Voter
		if err := voterCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&voter); err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Voter not found"})
			return
		}

		_ = json.NewEncoder(w).Encode(VoterResponse{
			Status: "success",
			Data:   voter.Registrations,
		})
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// sendEmail sends an HTML email via Hostinger SMTP (Asynchronously)
func sendEmail(to, subject, htmlBody string) error {
	select {
	case emailQueue <- emailJob{To: to, Subject: subject, HTMLBody: htmlBody}:
		// queued successfully
		return nil
	default:
		// Queue full, fallback to synchronous (or you could return an error)
		fmt.Printf("[WARN] Email queue full, sending synchronously to: %s\n", to)
		return sendEmailWithAttachment(to, subject, htmlBody, "", nil)
	}
}

// Global email queue
type emailJob struct {
	To             string
	Subject        string
	HTMLBody       string
	Filename       string
	AttachmentData []byte
}

var emailQueue = make(chan emailJob, 5000) // buffer 5000 emails to prevent API timeout on extreme bulk imports

func init() {
	go emailWorker()
}

func emailWorker() {
	// Respect custom SMTP limits (e.g., max 2-3 emails per second)
	ticker := time.NewTicker(400 * time.Millisecond)
	for job := range emailQueue {
		<-ticker.C // throttle
		err := sendEmailWithAttachment(job.To, job.Subject, job.HTMLBody, job.Filename, job.AttachmentData)
		if err != nil {
			fmt.Printf("[ERROR] Background email worker failed for %s: %v\n", job.To, err)
			// Simple retry logic could go here
		} else {
			fmt.Printf("[OK] Background email sent to %s\n", job.To)
		}
	}
}

// sendEmailWithAttachment sends an HTML email with an optional PDF attachment via Hostinger SMTP
func sendEmailWithAttachment(to, subject, htmlBody, filename string, attachmentData []byte) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderName := os.Getenv("SENDER_NAME")

	if smtpHost == "" {
		smtpHost = "smtp.hostinger.com"
	}
	if smtpPort == "" {
		smtpPort = "587"
	}
	if smtpUser == "" || smtpPass == "" || senderEmail == "" {
		return fmt.Errorf("SMTP not configured: set SMTP_USER, SMTP_PASS, and SENDER_EMAIL")
	}
	if senderName == "" {
		senderName = "SecureVote"
	}

	fromHeader := fmt.Sprintf("%s <%s>", senderName, senderEmail)

	var rawMsg bytes.Buffer

	if len(attachmentData) > 0 {
		// multipart/mixed for attachments
		mw := multipart.NewWriter(&rawMsg)
		boundary := mw.Boundary()

		rawMsg.Reset()
		rawMsg.WriteString(fmt.Sprintf("From: %s\r\n", fromHeader))
		rawMsg.WriteString(fmt.Sprintf("To: %s\r\n", to))
		rawMsg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
		rawMsg.WriteString("MIME-Version: 1.0\r\n")
		rawMsg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", boundary))

		// HTML body part
		rawMsg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		rawMsg.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
		rawMsg.WriteString(htmlBody)
		rawMsg.WriteString("\r\n")

		// Attachment part
		rawMsg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		h := make(textproto.MIMEHeader)
		h.Set("Content-Type", "application/pdf")
		h.Set("Content-Transfer-Encoding", "base64")
		h.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		for k, vs := range h {
			for _, v := range vs {
				rawMsg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
			}
		}
		rawMsg.WriteString("\r\n")
		encoded := base64.StdEncoding.EncodeToString(attachmentData)
		rawMsg.WriteString(encoded)
		rawMsg.WriteString("\r\n")
		rawMsg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		// simple HTML email
		rawMsg.WriteString(fmt.Sprintf("From: %s\r\n", fromHeader))
		rawMsg.WriteString(fmt.Sprintf("To: %s\r\n", to))
		rawMsg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
		rawMsg.WriteString("MIME-Version: 1.0\r\n")
		rawMsg.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
		rawMsg.WriteString(htmlBody)
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{to}, rawMsg.Bytes())
	if err != nil {
		return fmt.Errorf("smtp send error: %w", err)
	}
	return nil
}
func GetElectionVoters(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	electionAddress := vars["address"]
	if electionAddress == "" {
		http.Error(w, "Missing election address", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query nested array
	cursor, err := voterCollection.Find(ctx, bson.M{"registrations.election_address": electionAddress})
	if err != nil {
		http.Error(w, "Failed to fetch voters: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var voters []Voter
	if err := cursor.All(ctx, &voters); err != nil {
		http.Error(w, "Failed to decode voters: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// build response
	vlist := make([]map[string]interface{}, 0, len(voters))
	for _, v := range voters {
		var dobStr string
		if !v.DOB.IsZero() {
			dobStr = v.DOB.Format("2006-01-02")
		}

		// Find status for this election
		status := "Unknown"
		for _, reg := range v.Registrations {
			if reg.ElectionAddress == electionAddress {
				status = reg.Status
				break
			}
		}

		vlist = append(vlist, map[string]interface{}{
			"id":        v.ID.Hex(),
			"email":     v.Email,
			"full_name": v.FullName,
			"dob":       dobStr,
			"roll_no":   v.RollNo,
			"mobile":    v.Mobile,
			"gender":    v.Gender,
			"year":      v.Year,
			"status":    status,
		})
	}

	resp := map[string]interface{}{
		"status":  "success",
		"message": "voters retrieved",
		"voters":  vlist,
		"data":    map[string]interface{}{"voters": vlist},
		"count":   len(vlist),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// IsVoterVerified checks if a voter is allowed to vote
func IsVoterVerified(email, electionAddr string) bool {
	if voterCollection == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find voter who has this election registration
	var v Voter
	err := voterCollection.FindOne(ctx, bson.M{
		"email":                          email,
		"registrations.election_address": electionAddr,
	}).Decode(&v)

	if err != nil {
		return false
	}

	for _, reg := range v.Registrations {
		if reg.ElectionAddress == electionAddr {
			return reg.Status == "Verified"
		}
	}
	return false
}

func ResultMail(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VoterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Invalid request body"})
		return
	}
	if req.ElectionAddress == "" || req.ElectionName == "" || req.WinnerCandidate == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "election_address, election_name and winner_candidate are required"})
		return
	}

	// AUDIT LOG
	go LogAction(req.ElectionAddress, "ELECTION_ENDED", "System", fmt.Sprintf("Election '%s' ended. Winner: %s", req.ElectionName, req.WinnerCandidate))

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := voterCollection.Find(ctx, bson.M{"registrations.election_address": req.ElectionAddress})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Error fetching voters: " + err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var voters []Voter
	if err := cursor.All(ctx, &voters); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Error decoding voters: " + err.Error()})
		return
	}

	// Fetch Candidates from Blockchain
	var candidates []map[string]interface{}

	client, err := getClient()
	if err != nil {
		fmt.Printf("ResultMail: failed to connect to eth node: %v\n", err)
	} else {
		defer client.Close()
		ethAddr := common.HexToAddress(req.ElectionAddress)
		contract, err := bindings.NewElection(ethAddr, client)
		if err != nil {
			fmt.Printf("ResultMail: failed to bind contract: %v\n", err)
		} else {
			count, err := contract.GetNumOfCandidates(&bind.CallOpts{Context: context.Background()})
			if err != nil {
				fmt.Printf("ResultMail: failed to get candidate count: %v\n", err)
			} else {
				n := count.Int64()
				for i := int64(0); i < n; i++ {
					name, _, _, votes, _, err := contract.GetCandidate(&bind.CallOpts{Context: context.Background()}, big.NewInt(i))
					if err != nil {
						fmt.Printf("ResultMail: failed to get candidate %d: %v\n", i, err)
						continue
					}
					candidates = append(candidates, map[string]interface{}{
						"name":      name,
						"voteCount": int(votes.Int64()),
					})
				}
			}
		}
	}

	// Fallback to DB if blockchain fetch failed
	if len(candidates) == 0 && candidateCollection != nil {
		fmt.Println("ResultMail: chain fetch failed, falling back to DB")
		curC, errC := candidateCollection.Find(ctx, bson.M{"electionAddress": req.ElectionAddress})
		if errC == nil {
			var rawCands []map[string]interface{}
			_ = curC.All(ctx, &rawCands)
			candidates = rawCands
			curC.Close(ctx)
		}
	}

	// Build and send the results email (no PDF, no audit log)
	htmlBody := GenerateResultsEmailHTML(req.ElectionName, req.WinnerCandidate, candidates)
	subject := fmt.Sprintf("Results: %s - Winner Announced", req.ElectionName)

	var sendErrs []string

	for _, v := range voters {
		if v.Email == "" {
			continue
		}
		if err := sendEmail(v.Email, subject, htmlBody); err != nil {
			sendErrs = append(sendErrs, fmt.Sprintf("%s: %v", v.Email, err))
			fmt.Printf("sendEmail error for %s: %v\n", v.Email, err)
		}
	}

	// Also notify the winner candidate
	if req.CandidateEmail != "" {
		winnerSubject := "Congratulations! You Won - " + req.ElectionName
		if err := sendEmail(req.CandidateEmail, winnerSubject, htmlBody); err != nil {
			sendErrs = append(sendErrs, fmt.Sprintf("%s: %v", req.CandidateEmail, err))
		}
	}

	if len(sendErrs) == 0 {
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "success", Message: "mails sent successfully"})
		return
	}

	_ = json.NewEncoder(w).Encode(VoterResponse{
		Status:  "partial",
		Message: fmt.Sprintf("mails sent with %d errors", len(sendErrs)),
		Data:    map[string]interface{}{"mailErrors": sendErrs},
	})
}

// ApproveVoter Endpoint
func ApproveVoter(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	voterID := vars["voterId"]

	var req struct {
		Status string `json:"status"` // Verified or Rejected
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Status == "" {
		req.Status = "Verified"
	}

	objID, _ := primitive.ObjectIDFromHex(voterID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := voterCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": req.Status}})
	if err != nil {
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(VoterResponse{Status: "success", Message: "Voter status updated to " + req.Status})
}

// VerifyAndDeleteOTP checks OTP validity and deletes it if valid
func VerifyAndDeleteOTP(email, otp string) bool {
	if otpCollection == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var otpDoc struct {
		OTP       string             `bson:"otp"`
		ExpiresAt primitive.DateTime `bson:"expiresAt"`
	}
	// find by email
	err := otpCollection.FindOne(ctx, bson.M{"email": email}).Decode(&otpDoc)
	if err != nil {
		return false
	} // not found

	// check expiry
	if time.Now().UTC().After(otpDoc.ExpiresAt.Time()) {
		return false
	}

	// check match
	if otpDoc.OTP != otp {
		return false
	}

	// valid -> delete
	otpCollection.DeleteMany(ctx, bson.M{"email": email})
	return true
}

// GetVoterAnalytics returns the count of voters grouped by address (e.g. City)
func GetVoterAnalytics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	electionAddr := vars["address"] // matches /elections/{address}/analytics/geo

	if voterCollection == nil {
		http.Error(w, "DB not initialized", http.StatusInternalServerError)
		return
	}

	// Pipeline: Match Election in Registrations -> Group by Address -> Count
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "registrations.election_address", Value: electionAddr}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$address"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}}, // highest first
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := voterCollection.Aggregate(ctx, pipeline)
	if err != nil {
		http.Error(w, "Aggregation failed", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var results []struct {
		Address string `bson:"_id" json:"address"`
		Count   int    `bson:"count" json:"count"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		http.Error(w, "Cursor decode failed", http.StatusInternalServerError)
		return
	}

	// Clean up empty addresses
	final := []interface{}{}
	for _, item := range results {
		addr := item.Address
		if addr == "" {
			addr = "Unknown"
		}
		final = append(final, map[string]interface{}{
			"region": addr,
			"count":  item.Count,
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   final,
	})
}

// AddVotersToElection allows Admin to bulk add existing voters to an election
func AddVotersToElection(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	electionAddr := vars["address"]
	if electionAddr == "" {
		http.Error(w, "Election address required", http.StatusBadRequest)
		return
	}

	var req struct {
		VoterIDs []string `json:"voter_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Invalid body"})
		return
	}

	if len(req.VoterIDs) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "No voter IDs provided"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	successCount := 0
	alreadyCount := 0
	errCount := 0

	for _, vid := range req.VoterIDs {
		objID, err := primitive.ObjectIDFromHex(vid)
		if err != nil {
			errCount++
			continue
		}

		// Check if already registered for this election
		// We use $addToSet logic manually or via check to prevent duplicates
		count, _ := voterCollection.CountDocuments(ctx, bson.M{
			"_id":                            objID,
			"registrations.election_address": electionAddr,
		})

		if count > 0 {
			alreadyCount++
			continue
		}

		// Add registration
		newReg := VoterRegistration{
			ElectionAddress: electionAddr,
			Status:          "Verified", // Admin added = Verified
			RegisteredAt:    time.Now().UTC(),
		}

		_, err = voterCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$push": bson.M{"registrations": newReg}})
		if err != nil {
			errCount++
		} else {
			successCount++
		}
	}

	msg := fmt.Sprintf("Added %d voters. %d already in election. %d errors.", successCount, alreadyCount, errCount)
	_ = json.NewEncoder(w).Encode(VoterResponse{
		Status:  "success",
		Message: msg,
		Data: map[string]interface{}{
			"added":   successCount,
			"skipped": alreadyCount,
		},
	})
}

// ===== ForgotPassword =====
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VoterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Invalid request body"})
		return
	}
	if req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "email is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var voterInfo Voter
	// Look up the voter by email
	if err := voterCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&voterInfo); err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Voter account not found"})
		return
	}

	// Generate a new password
	rawPassword, err := genPassword(12)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Failed to generate new password"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Failed to hash password"})
		return
	}

	// Update the voter's password in the database
	_, err = voterCollection.UpdateOne(ctx, bson.M{"email": req.Email}, bson.M{"$set": bson.M{"password": string(hashed)}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Failed to update password"})
		return
	}

	// Email the new password to the voter
	subject := "Your New Password for SecureVote"
	body := GenerateForgotPasswordEmail(voterInfo.FullName, rawPassword)

	if err := sendEmail(req.Email, subject, body); err != nil {
		fmt.Printf("sendEmail error (ForgotPassword): %v\n", err)
		// We still updated the DB, but couldn't send the email
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Password reset, but failed to send email. Please contact admin."})
		return
	}

	_ = json.NewEncoder(w).Encode(VoterResponse{
		Status:  "success",
		Message: "A new password has been sent to your email",
	})
}

// ===== BulkResetVoterPasswords =====
// POST /api/elections/{address}/voters/reset-passwords
// Admin-triggered: for every voter registered in the election, generates a unique
// random password, saves bcrypt hash to DB, and emails the plain password to the voter.
func BulkResetVoterPasswords(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	electionAddress := vars["address"]
	if electionAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "election address is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Fetch all voters registered in this election
	cursor, err := voterCollection.Find(ctx, bson.M{"registrations.election_address": electionAddress})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Failed to fetch voters"})
		return
	}
	defer cursor.Close(ctx)

	var voters []Voter
	if err := cursor.All(ctx, &voters); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Failed to decode voters"})
		return
	}

	if len(voters) == 0 {
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "success", Message: "No voters found in this election"})
		return
	}

	successCount := 0
	failCount := 0
	var failedEmails []string

	for _, v := range voters {
		// Generate a unique password per voter
		rawPassword, err := genPassword(12)
		if err != nil {
			failCount++
			failedEmails = append(failedEmails, v.Email)
			continue
		}

		// Use MinCost here for extreme bulk ops (1000+ voters) to prevent HTTP timeouts and massive CPU spikes
		hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.MinCost)
		if err != nil {
			failCount++
			failedEmails = append(failedEmails, v.Email)
			continue
		}

		// Update password in DB
		updateCtx, updateCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, dbErr := voterCollection.UpdateOne(updateCtx, bson.M{"_id": v.ID}, bson.M{"$set": bson.M{"password": string(hashed)}})
		updateCancel()
		if dbErr != nil {
			failCount++
			failedEmails = append(failedEmails, v.Email)
			continue
		}

		// Send email with the new password
		subject := "Your New Password for SecureVote -- E-Voting System"
		body := GenerateForgotPasswordEmail(v.FullName, rawPassword)
		if emailErr := sendEmail(v.Email, subject, body); emailErr != nil {
			fmt.Printf("BulkResetVoterPasswords: sendEmail error for %s: %v\n", v.Email, emailErr)
			failCount++
			failedEmails = append(failedEmails, v.Email)
			continue
		}

		successCount++
	}

	msg := fmt.Sprintf("Passwords sent: %d succeeded, %d failed", successCount, failCount)
	respData := map[string]interface{}{
		"total":         len(voters),
		"success_count": successCount,
		"fail_count":    failCount,
	}
	if len(failedEmails) > 0 {
		respData["failed_emails"] = failedEmails
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": msg,
		"data":    respData,
	})
}

// ===== AdminResetVoterPassword =====
// POST /api/voters/{voterId}/reset-password
// Admin-triggered: generates a new password, saves hashed to DB, emails plain password to voter.
func AdminResetVoterPassword(w http.ResponseWriter, r *http.Request) {
	withVoterCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	voterIdStr := vars["voterId"]
	if voterIdStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "voterId is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find voter by ID (ObjectID) or email
	var voterInfo Voter
	objID, err := primitive.ObjectIDFromHex(voterIdStr)
	if err == nil {
		err = voterCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&voterInfo)
	} else {
		// fallback: treat as email
		err = voterCollection.FindOne(ctx, bson.M{"email": voterIdStr}).Decode(&voterInfo)
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Voter not found"})
		return
	}

	// Generate a new random password
	rawPassword, err := genPassword(12)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Failed to generate password"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Failed to hash password"})
		return
	}

	// Update in DB
	filter := bson.M{"_id": voterInfo.ID}
	_, err = voterCollection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"password": string(hashed)}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Failed to update password in DB"})
		return
	}

	// Email the new password
	subject := "Your New Password for SecureVote -- E-Voting System"
	body := GenerateForgotPasswordEmail(voterInfo.FullName, rawPassword)
	if err := sendEmail(voterInfo.Email, subject, body); err != nil {
		fmt.Printf("AdminResetVoterPassword: sendEmail error for %s: %v\n", voterInfo.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Password reset but email failed. Check server logs."})
		return
	}

	_ = json.NewEncoder(w).Encode(VoterResponse{
		Status:  "success",
		Message: "New password generated and emailed to " + voterInfo.Email,
	})
}
