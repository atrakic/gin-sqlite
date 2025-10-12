package auth

import (
	"errors"
	"os"
	"time"

	"github.com/atrakic/gin-sqlite/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var (
	// Default JWT secret key - should be set via environment variable in production
	jwtSecret = []byte("your-secret-key")
)

func init() {
	// Override with environment variable if set
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		jwtSecret = []byte(secret)
	}
}

// GenerateJWT generates a JWT token for the given username
func GenerateJWT(username string) (string, int64, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &models.JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gin-sqlite-demo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationTime.Unix(), nil
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateCredentials validates username and password
func ValidateCredentials(username, password string) bool {
	// Get credentials from environment variables
	adminUser := os.Getenv("ADMIN_USER")
	if adminUser == "" {
		adminUser = "admin"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "secret"
	}

	return username == adminUser && password == adminPassword
}
