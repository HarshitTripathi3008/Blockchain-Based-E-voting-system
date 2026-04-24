package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey    contextKey = "userID"
	UserEmailKey contextKey = "userEmail"
	UserRoleKey  contextKey = "userRole"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "securevote-fallback-secret-key-123" // Fallback, override in prod
	}
	return []byte(secret)
}

// verifyToken checks for a token in the Authorization header OR a specific cookie
func verifyToken(r *http.Request, allowedRole, cookieName string) (*jwt.Token, *jwt.MapClaims, error) {
	var tokenStr string

	// 1. Try from Header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 2. Try from Cookie (for seamless frontend transition)
	if tokenStr == "" && cookieName != "" {
		if cookie, err := r.Cookie(cookieName); err == nil {
			tokenStr = cookie.Value
		}
	}

	if tokenStr == "" {
		return nil, nil, http.ErrNoCookie
	}

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		return nil, nil, err
	}

	// verify role if specified
	if allowedRole != "" {
		role, ok := (*claims)["role"].(string)
		if !ok || role != allowedRole {
			return nil, nil, jwt.ErrInvalidKey
		}
	}

	return token, claims, nil
}

// VoterAuthMiddleware protects routes meant for verified voters
func VoterAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow options bypass for CORS preflights
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		_, claims, err := verifyToken(r, "voter", "voter_token")
		if err != nil {
			http.Error(w, `{"status":"error", "message":"Unauthorized: Invalid or missing token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, (*claims)["sub"])
		ctx = context.WithValue(ctx, UserEmailKey, (*claims)["email"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CompanyAuthMiddleware protects admin/company routes
func CompanyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow options bypass for CORS preflights
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		_, claims, err := verifyToken(r, "company", "company_token")
		if err != nil {
			http.Error(w, `{"status":"error", "message":"Unauthorized: Admin access required"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, (*claims)["sub"])
		ctx = context.WithValue(ctx, UserEmailKey, (*claims)["email"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
