package utils //nolint:typecheck

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"os"
	"strconv"
	"time"
)

var (
	ErrSecretKeyIsNull      = errors.New("env varriable SECRET_KEY is not set or is empty.")
	ErrMalformedToken       = errors.New("malformed token")
	ErrorExpiredOrNotActive = errors.New("token is either expired or not active yet")
)

type Claims struct {
	UserID uuid.UUID `json:"userID"`
	jwt.StandardClaims
}

type User struct {
	UserID uuid.UUID
}

func GenerateJWT(user User) (string, error) {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return "", ErrSecretKeyIsNull
	}
	jwtExpirationTime := os.Getenv("JWT_EXPIRATION_TIME")

	intExpirationTime := 12

	if jwtExpirationTime != "" {
		t, err := strconv.Atoi(jwtExpirationTime)
		if err == nil {
			intExpirationTime = t
		}
	}

	var jwtKey = []byte(secret)

	expirationTime := time.Now().Add(time.Duration(intExpirationTime) * time.Hour)

	claims := &Claims{
		UserID: user.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(signedToken string) (*Claims, error) {
	claims := &Claims{}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return nil, ErrSecretKeyIsNull
	}

	var jwtKey = []byte(secret)
	token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, err
		}

		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrMalformedToken
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, ErrorExpiredOrNotActive
			}
			return nil, err
		}
	}

	if !token.Valid {
		return nil, err
	}
	return claims, nil
}
