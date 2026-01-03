package utils

import (
	"advancedBackend/configs"
	"advancedBackend/errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

func GenerateToken(userID string, email string) (string, error) {
	// Generating JWT token for a user using jwt

	serverConfig := configs.GetServerConfig()
	secretKeyString := serverConfig.SecretKey
	secretKey := []byte(secretKeyString)

	// Generating the token using the secret key
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userID,
			"email":  email,
			"exp":    time.Now().Add(time.Hour + 24).Unix(),
		},
	)
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		slog.Error(
			"Error while signing the token with secret key",
			slog.Any("Error", err),
		)
		_, customError := errors.NewInternalServerError("Error while signing the token with secret key", err)
		return "", customError
	}

	return tokenString, nil
}

func ValidateToken(token string) (*Claims, error, []byte, int) {
	// Validating the JWT token using the secret key

	serverConfig := configs.GetServerConfig()
	secretKeyString := serverConfig.SecretKey
	secretKey := []byte(secretKeyString)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		slog.Error(
			"Error while parsing the token",
			slog.Any("Error", err),
		)
		errorJson, internalServerError := errors.NewInternalServerError("Error while parsing the token", err)
		return nil, internalServerError, errorJson, internalServerError.Code
	}

	if !parsedToken.Valid {
		slog.Error(
			"Invalid Token passed",
			slog.Any("Error", err),
		)
		errorJson, badRequestError := errors.NewBadRequestError("Invalid Token passed", err)
		return nil, badRequestError, errorJson, badRequestError.Code
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		slog.Error(
			"Error while converting the claims",
			slog.Any("Error", err),
		)
		errorJson, internalServerError := errors.NewInternalServerError("Error while converting the claims", err)
		return nil, internalServerError, errorJson, internalServerError.Code
	}

	userID := claims["userId"].(string)
	email := claims["email"].(string)

	return &Claims{
		UserID: userID,
		Email:  email,
	}, nil, nil, 0
}
