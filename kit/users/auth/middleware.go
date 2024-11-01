package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/juaguz/yuno/kit/users/dto"
)

type UserContextKey string

const UserKey UserContextKey = "user"

type UserRepository interface {
	FindByExternalID(ctx context.Context, externalID string) (*dto.User, error)
}

type UserClaims struct {
	Username string `json:"preferred_username"`
	Email    string `json:"email"`
	UserID   string `json:"sub"`
	jwt.RegisteredClaims
}

// Fetches the public key from Keycloak based on the given kid
func getKeycloakPublicKey(kid, keycloakURL, realm string) (*rsa.PublicKey, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", keycloakURL, realm)

	// Make HTTP GET request to fetch the JSON of public keys
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public keys: %v", err)
	}
	defer resp.Body.Close()

	var keyResponse struct {
		Keys []struct {
			Kid     string   `json:"kid"`
			Kty     string   `json:"kty"`
			Alg     string   `json:"alg"`
			Use     string   `json:"use"`
			N       string   `json:"n"`
			E       string   `json:"e"`
			X5c     []string `json:"x5c"`
			X5t     string   `json:"x5t"`
			X5tS256 string   `json:"x5t#S256"`
		} `json:"keys"`
	}
	b, _ := io.ReadAll(resp.Body)
	json.Unmarshal(b, &keyResponse)

	// Find the correct key by kid
	for _, key := range keyResponse.Keys {
		if key.Kid == kid {
			// Decode the modulus and exponent from base64
			nb, _ := base64.RawURLEncoding.DecodeString(key.N)
			eb, _ := base64.RawURLEncoding.DecodeString(key.E)

			// Convert exponent bytes to integer
			e := int(new(big.Int).SetBytes(eb).Uint64())

			// Construct the RSA public key
			pubKey := &rsa.PublicKey{
				N: new(big.Int).SetBytes(nb),
				E: e,
			}
			return pubKey, nil
		}
	}

	return nil, fmt.Errorf("public key with kid %s not found", kid)
}

func JWTMiddleware(userRepo UserRepository, keycloakURL, realm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip JWT validation for Swagger routes
			if strings.HasPrefix(r.URL.Path, "/swagger/") {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Malformed Authorization header", http.StatusUnauthorized)
				return
			}

			// Parse the token without verifying to get the kid
			token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
				kid, ok := token.Header["kid"].(string)
				if !ok {
					return nil, fmt.Errorf("kid not found in token header")
				}

				return getKeycloakPublicKey(kid, keycloakURL, realm)
			})
			if err != nil || !token.Valid {

				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Extract and validate claims
			if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
				// Find the user by external ID in the database
				u, err := userRepo.FindByExternalID(r.Context(), claims.UserID)
				if err != nil {
					http.Error(w, "User not found", http.StatusUnauthorized)
					return
				}
				// Store user in context
				ctx := context.WithValue(r.Context(), UserKey, u)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			}
		})
	}
}

func GetUserFromContext(ctx context.Context) (*dto.User, error) {
	u, ok := ctx.Value(UserKey).(*dto.User)
	if !ok {
		return nil, fmt.Errorf("no user claims in context")
	}
	return u, nil
}
