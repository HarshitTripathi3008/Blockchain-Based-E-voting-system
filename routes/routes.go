package routes

import (
	"log"
	"net/http"

	"MAJOR-PROJECT/controllers"
	"MAJOR-PROJECT/middleware"

	"github.com/gorilla/mux"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/").Subrouter()
	api.Use(middleware.RecoveryMiddleware) // Global panic recovery
	api.Use(corsMiddleware)
	api.Use(middleware.RateLimitMiddleware)

	// ----------------------------
	// COMPANY ROUTES
	// ----------------------------
	api.HandleFunc("/admin/clear-database", controllers.ClearDatabase).Methods(http.MethodPost, http.MethodOptions) // NEW
	api.HandleFunc("/company/register", controllers.CreateCompany).Methods(http.MethodPost, http.MethodOptions)

	api.HandleFunc("/company/authenticate", controllers.AuthenticateCompany).Methods(http.MethodPost, http.MethodOptions)

	// ----------------------------
	// ELECTION ROUTES
	// ----------------------------
	api.HandleFunc("/elections/create", controllers.CreateElection).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/elections/{address}/details", controllers.GetElectionInfo).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/elections/{address}/candidates", controllers.GetElectionCandidates).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/elections/{address}/vote", controllers.VoteCandidate).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/elections/{address}/voters", controllers.GetElectionVoters).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/elections/dates", controllers.SetElectionDates).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/elections/{address}/metadata", controllers.GetElectionMetadata).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/elections/{address}/analytics/geo", controllers.GetVoterAnalytics).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/elections/{address}/end", controllers.EndElection).Methods(http.MethodPost, http.MethodOptions)  // NEW
	api.HandleFunc("/elections", controllers.GetAllElections).Methods(http.MethodGet, http.MethodOptions)             // NEW
	api.HandleFunc("/elections/archives", controllers.GetArchivedResults).Methods(http.MethodGet, http.MethodOptions) // L1 Archives

	// ----------------------------
	// CANDIDATE ROUTES
	// ----------------------------
	api.HandleFunc("/candidate/register", controllers.RegisterCandidate).Methods(http.MethodPost, http.MethodOptions)

	// ----------------------------
	// VOTER ROUTES
	// ----------------------------
	api.HandleFunc("/voters/register", controllers.RegisterVoter).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/voters/send-otp", controllers.SendOTP).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/voters/verify-otp-register", controllers.VerifyOTPAndRegister).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/voters/me/elections", controllers.GetVoterElections).Methods(http.MethodGet, http.MethodOptions) // NEW
	api.HandleFunc("/voters/forgot-password", controllers.ForgotPassword).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/voter/authenticate", controllers.AuthenticateVoter).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/voters", controllers.GetAllVoters).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	api.HandleFunc("/voters/{voterId}/card", controllers.GenerateVoterID).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/voters/{voterId}/card/email", controllers.EmailVoterID).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/voters/{voterId}", controllers.UpdateVoter).Methods(http.MethodPut, http.MethodOptions)
	api.HandleFunc("/voters/{voterId}", controllers.DeleteVoter).Methods(http.MethodDelete, http.MethodOptions)
	api.HandleFunc("/voters/{voterId}/approve", controllers.ApproveVoter).Methods(http.MethodPost, http.MethodOptions)                              // NEW
	api.HandleFunc("/voters/{voterId}/reset-password", controllers.AdminResetVoterPassword).Methods(http.MethodPost, http.MethodOptions)            // ADMIN RESET
	api.HandleFunc("/elections/{address}/voters/add", controllers.AddVotersToElection).Methods(http.MethodPost, http.MethodOptions)                 // NEW BULK IMPORT
	api.HandleFunc("/elections/{address}/voters/reset-passwords", controllers.BulkResetVoterPasswords).Methods(http.MethodPost, http.MethodOptions) // BULK SEND PASSWORDS
	api.HandleFunc("/voter/resultMail", controllers.ResultMail).Methods(http.MethodPost, http.MethodOptions)
	// ----------------------------
	// UPLOAD ROUTES
	// ----------------------------
	// Unified hybrid upload route
	api.HandleFunc("/upload/unified", controllers.UnifiedUploadHandler).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/upload/s3", controllers.UnifiedUploadHandler).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/upload/gdrive", controllers.UnifiedUploadHandler).Methods(http.MethodPost, http.MethodOptions)

	// ----------------------------
	// STATIC FILE SERVING
	// ----------------------------
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Root handler serves homepage.html
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pages/homepage.html")
	})

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./pages/")))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 - %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})

	// Favicon handler to avoid 404 noise
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	return router
}
