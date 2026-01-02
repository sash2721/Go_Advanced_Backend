package services

import (
	"advancedBackend/errors"
	"advancedBackend/utils"
	"log/slog"
	"strings"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// imitating the repository layer for now
var passwordCache = map[string]string{
	"sakshipaygude27@gmail.com": "",
	"sahilshah2104@gmail.com":   "",
}

func (*AuthService) Login(email string, password string) (string, error, []byte) {
	// get the hashed password from the passwordCache
	hashedPassword := passwordCache[email]

	// send this hashedPassword and original password for comparison
	err := utils.ComparePassword(hashedPassword, password)

	if err != nil {
		slog.Error(
			"Password is incorrect, please retry",
			slog.Any("Error", err),
		)
		errorJson, badRequestError := errors.NewBadRequestError("Password is incorrect, please retry", err)
		return "", badRequestError.Error, errorJson
	}

	// extracting the userID from the mail itself
	userID := strings.Split(email, "@")[0]

	// generating the JWT token for this part
	jwtToken, err := utils.GenerateToken(userID, email)

	if err != nil {
		slog.Error(
			"Error while generating the JWT token",
			slog.Any("Error", err),
		)
		errorJson, internalServerError := errors.NewInternalServerError("Error while generating the JWT token", err)
		return "", internalServerError.Error, errorJson
	}

	return jwtToken, nil, nil
}
