package auth

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// Claims struct to be encoded to JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// GenerateToken generates a JWT token for a given username
func GenerateToken(username string) (string, error) {
	// The expiration time after which the token will be invalid.
	expirationTime := time.Now().Add(24 * time.Hour).Unix()

	// Create the JWT claims, which includes the username and expiration time
	claims := &jwt.StandardClaims{
		// In JWT, the expiry time is expressed as unix milliseconds
		ExpiresAt: expirationTime,
		Issuer:    username,
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(JwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRandomKey generates a random key for JWT signing
func GenerateRandomKey() string {
	key := make([]byte, 32) // generate a 256 bit key
	_, err := rand.Read(key)
	if err != nil {
		panic("Failed to generate random key: " + err.Error())
	}

	return base64.StdEncoding.EncodeToString(key)
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ComparePassword(dbPassword string, incomingPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(incomingPassword))
	if err != nil {
		return err
	}

	return nil
}
