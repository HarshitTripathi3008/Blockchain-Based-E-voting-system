package controllers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Voter struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email           string             `bson:"email" json:"email"`
	Password        string             `bson:"password" json:"-"`
	FullName        string             `bson:"full_name,omitempty" json:"full_name,omitempty"`
	DOB             time.Time          `bson:"dob,omitempty" json:"dob,omitempty"`
	Address         string             `bson:"address,omitempty" json:"address,omitempty"`
	Mobile          string             `bson:"mobile,omitempty" json:"mobile,omitempty"`
	FatherName      string             `bson:"father_name,omitempty" json:"father_name,omitempty"`
	MotherName      string             `bson:"mother_name,omitempty" json:"mother_name,omitempty"`
	ElectionAddress string             `bson:"election_address" json:"election_address"`
}

type VoterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	FullName        string `json:"full_name,omitempty"`
	DOB             string `json:"dob,omitempty"` // expect "YYYY-MM-DD"
	Address         string `json:"address,omitempty"`
	Mobile          string `json:"mobile,omitempty"`
	FatherName      string `json:"father_name,omitempty"`
	MotherName      string `json:"mother_name,omitempty"`
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

var voterCollection *mongo.Collection
var otpCollection *mongo.Collection

// Initialize collections
func InitVoterCollection(client *mongo.Client, dbName string) {
	voterCollection = client.Database(dbName).Collection("voters")
	fmt.Println("✅ Initialized voters collection")
}

func InitOTPCollection(client *mongo.Client, dbName string) {
	otpCollection = client.Database(dbName).Collection("otps")
	fmt.Println("✅ Initialized OTP collection")
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
	body := fmt.Sprintf(`<p>Your OTP for voter registration is <b>%s</b>. It expires in 10 minutes.</p><p>If you did not request this, ignore this email.</p>`, otp)
	if err := sendEmail(req.Email, subject, body); err != nil {
		fmt.Printf("sendEmail error (SendOTP): %v\n", err)
		// remove inserted OTP because email failed
		_, _ = otpCollection.DeleteMany(ctx, bson.M{"email": req.Email})
		http.Error(w, "failed to send otp email", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(VoterResponse{Status: "success", Message: "OTP sent"})
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
	if req.Email == "" || req.OTP == "" {
		http.Error(w, "email and otp required", http.StatusBadRequest)
		return
	}
	if otpCollection == nil || voterCollection == nil {
		http.Error(w, "server misconfigured: collections not ready", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// fetch OTP doc
	var otpDoc struct {
		OTP       string             `bson:"otp"`
		ExpiresAt primitive.DateTime `bson:"expiresAt"`
	}
	if err := otpCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&otpDoc); err != nil {
		http.Error(w, "OTP not found or expired", http.StatusBadRequest)
		return
	}
	if time.Now().UTC().After(otpDoc.ExpiresAt.Time()) || req.OTP != otpDoc.OTP {
		http.Error(w, "invalid or expired otp", http.StatusBadRequest)
		return
	}

	// basic validations
	if req.DOB != "" {
		if dob, err := parseDOB(req.DOB); err == nil {
			if computeAge(dob) < 18 {
				http.Error(w, "voter must be 18+", http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "invalid dob", http.StatusBadRequest)
			return
		}
	}
	if req.Mobile != "" && !mobileRe.MatchString(req.Mobile) {
		http.Error(w, "invalid mobile format", http.StatusBadRequest)
		return
	}

	// ensure email not already registered for same election
	count, err := voterCollection.CountDocuments(ctx, bson.M{
		"email":            req.Email,
		"election_address": req.ElectionAddress,
	})
	if err == nil && count > 0 {
		http.Error(w, "voter with this email already registered for this election", http.StatusConflict)
		return
	}

	// generate password
	rawPassword, err := genPassword(12)
	if err != nil {
		http.Error(w, "failed to create password", http.StatusInternalServerError)
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	dobTime, _ := parseDOB(req.DOB)

	newVoter := Voter{
		Email:           req.Email,
		Password:        string(hashed),
		FullName:        req.FullName,
		DOB:             dobTime,
		Mobile:          req.Mobile,
		Address:         req.Address,
		FatherName:      req.FatherName,
		MotherName:      req.MotherName,
		ElectionAddress: req.ElectionAddress,
	}

	if _, err := voterCollection.InsertOne(ctx, newVoter); err != nil {
		http.Error(w, "failed to create voter", http.StatusInternalServerError)
		return
	}

	// delete OTP record
	_, _ = otpCollection.DeleteMany(ctx, bson.M{"email": req.Email})

	// email password to voter (plaintext in email — consider password reset flow for production)
	subject := "Your voter account credentials"
	body := fmt.Sprintf(`<p>Hello %s,</p><p>Your voter account has been created for election <b>%s</b>.</p><p><b>Email:</b> %s<br/><b>Password:</b> %s</p><p>Please login and change your password.</p>`, req.FullName, req.ElectionAddress, req.Email, rawPassword)
	if err := sendEmail(req.Email, subject, body); err != nil {
		fmt.Printf("sendEmail error (password email): %v\n", err)
		// do not fail the request — account created
	}

	_ = json.NewEncoder(w).Encode(VoterResponse{Status: "success", Message: "voter created; password emailed"})
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

	// parse and validate DOB → age >= 18
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

	// default password = email if none provided (existing behavior preserved)
	rawPassword := req.Password
	if rawPassword == "" {
		rawPassword = req.Email
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Error hashing password"})
		return
	}

	newVoter := Voter{
		Email:           req.Email,
		Password:        string(hashed),
		FullName:        req.FullName,
		DOB:             dob,
		Address:         req.Address,
		Mobile:          req.Mobile,
		FatherName:      req.FatherName,
		MotherName:      req.MotherName,
		ElectionAddress: req.ElectionAddress,
	}

	result, err := voterCollection.InsertOne(ctx, newVoter)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Code == 11000 {
					w.WriteHeader(http.StatusConflict)
					_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Voter already exists"})
					return
				}
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "Voter could not be added"})
		return
	}

	var created Voter
	_ = voterCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&created)

	// Return minimal info (no password) + created fields
	_ = json.NewEncoder(w).Encode(VoterResponse{
		Status:  "success",
		Message: "Voter added successfully.",
		Data: map[string]interface{}{
			"id":            created.ID.Hex(),
			"email":         created.Email,
			"full_name":     created.FullName,
			"dob":           created.DOB.Format("2006-01-02"),
			"address":       created.Address,
			"mobile":        created.Mobile,
			"father_name":   created.FatherName,
			"mother_name":   created.MotherName,
			"election_addr": created.ElectionAddress,
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
			"id":               voterInfo.ID.Hex(),
			"election_address": voterInfo.ElectionAddress,
			"full_name":        voterInfo.FullName,
			"mobile":           voterInfo.Mobile,
		},
	})
}

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
	if electionAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "election_address is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := voterCollection.Find(ctx, bson.M{"election_address": electionAddress})
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
		vlist = append(vlist, map[string]interface{}{
			"id":          id,
			"email":       v.Email,
			"full_name":   v.FullName,
			"dob":         dobStr,
			"address":     v.Address,
			"mobile":      v.Mobile,
			"father_name": v.FatherName,
			"mother_name": v.MotherName,
		})
	}

	resp := map[string]interface{}{
		"status":  "success",
		"message": "voters list found",
		"voters":  vlist,
		"data":    map[string]interface{}{"voters": vlist},
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
	if req.Address != "" {
		update["address"] = req.Address
	}
	if req.Mobile != "" {
		if !mobileRe.MatchString(req.Mobile) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(VoterResponse{Status: "error", Message: "invalid mobile format"})
			return
		}
		update["mobile"] = req.Mobile
	}
	if req.FatherName != "" {
		update["father_name"] = req.FatherName
	}
	if req.MotherName != "" {
		update["mother_name"] = req.MotherName
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

// sendEmail using SendGrid API
func sendEmail(to, subject, htmlBody string) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	fromEmail := os.Getenv("SENDER_EMAIL")

	if apiKey == "" || fromEmail == "" {
		return fmt.Errorf("sendgrid not configured: set SENDGRID_API_KEY and SENDER_EMAIL")
	}

	from := mail.NewEmail("Voting System", fromEmail)
	toAddr := mail.NewEmail("", to)

	// Create message; using htmlBody as HTML content and leaving text content empty
	message := mail.NewSingleEmail(from, subject, toAddr, "", htmlBody)

	client := sendgrid.NewSendClient(apiKey)
	resp, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("sendgrid send error: %w", err)
	}
	// SendGrid returns 202 Accepted on success
	if resp.StatusCode >= 300 {
		return fmt.Errorf("sendgrid returned status %d: %s", resp.StatusCode, resp.Body)
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

	cursor, err := voterCollection.Find(ctx, bson.M{"election_address": electionAddress})
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

	// build response with expanded fields (same shape as GetAllVoters)
	vlist := make([]map[string]interface{}, 0, len(voters))
	for _, v := range voters {
		var dobStr string
		if !v.DOB.IsZero() {
			dobStr = v.DOB.Format("2006-01-02")
		}
		vlist = append(vlist, map[string]interface{}{
			"id":          v.ID.Hex(),
			"email":       v.Email,
			"full_name":   v.FullName,
			"dob":         dobStr,
			"address":     v.Address,
			"mobile":      v.Mobile,
			"father_name": v.FatherName,
			"mother_name": v.MotherName,
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := voterCollection.Find(ctx, bson.M{"election_address": req.ElectionAddress})
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

	subject := req.ElectionName + " - Election Results"
	body := fmt.Sprintf(`<h2>Election Results</h2><p>The results of <b>%s</b> are out.</p><p><b>Winner:</b> %s</p><br><p>Thank you for participating.</p>`, req.ElectionName, req.WinnerCandidate)

	var sendErrs []string

	for _, v := range voters {
		if v.Email == "" {
			continue
		}
		if err := sendEmail(v.Email, subject, body); err != nil {
			sendErrs = append(sendErrs, fmt.Sprintf("%s: %v", v.Email, err))
			fmt.Printf("sendEmail error for %s: %v\n", v.Email, err)
		}
	}

	// winner mail
	if req.CandidateEmail != "" {
		winnerBody := fmt.Sprintf(`<h2>Congratulations!</h2><p>You have won the <b>%s</b> election.</p><br><p>Best regards,<br/>Voting System</p>`, req.ElectionName)
		if err := sendEmail(req.CandidateEmail, subject, winnerBody); err != nil {
			sendErrs = append(sendErrs, fmt.Sprintf("%s: %v", req.CandidateEmail, err))
			fmt.Printf("sendEmail error for winner %s: %v\n", req.CandidateEmail, err)
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


