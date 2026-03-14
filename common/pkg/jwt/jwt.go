package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	IdentityID   uint   `json:"identity_id"`
	IdentityType string `json:"identity_type"`
	ClientID     *uint  `json:"client_id,omitempty"`
	EmployeeID   *uint  `json:"employee_id,omitempty"`
	jwt.RegisteredClaims
}

// JWTVerifier validates tokens locally using a shared HMAC secret.
// Every service creates one of these with the same secret that
// user-service uses to sign tokens.
type JWTVerifier struct {
	secret []byte
}

func NewJWTVerifier(secret string) *JWTVerifier {
	return &JWTVerifier{
		secret: []byte(secret),
	}
}

func GenerateToken(claims *Claims, secret string, expiryMinutes int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)

	tokenClaims := *claims
	tokenClaims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (v *JWTVerifier) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return v.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}
