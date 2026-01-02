package middlewares

import (
	"advancedBackend/utils"
	"context"
	"log/slog"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract the JWT auth token from the headers
		authzToken := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authzToken, "Bearer ")

		// validate the auth token
		userInfo, err, errJson, errorCode := utils.ValidateToken(tokenString)

		if err != nil {
			if errorCode == http.StatusInternalServerError {
				slog.Error(
					"Some internal issue occured while validating the token",
					slog.Any("Error", err),
				)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(errJson)
				return
			} else if errorCode == http.StatusBadRequest {
				slog.Debug(
					"Invalid Token sent in the request",
					slog.Any("Error", err),
				)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(errJson)
				return
			}
		}

		slog.Info(
			"Successfully Validated the JWT Token",
		)

		// Extracting the user info
		userId := userInfo.UserID
		email := userInfo.Email

		// storing the user info in the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userId)
		ctx = context.WithValue(ctx, "email", email)

		slog.Info(
			"Stored the User Info in the request context",
			slog.String("UserID", userId),
			slog.String("Email", email),
		)

		// Call the next handler in the line
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
