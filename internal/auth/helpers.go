package auth

import (
	"context"
	"errors"
	"net/http"
	"svipp-server/internal/database"
)

func CreateCookie(token string) http.Cookie {
	return http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   sessionExpirationSeconds, // Set the expiration time in seconds (e.g., 1 hour)
	}
}

func ClearCookie() *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

func GetSessionFromCtx(ctx context.Context) (database.GetSessionRow, error) {
	session, ok := ctx.Value(SessionContextKey).(database.GetSessionRow)
	if !ok {
		return database.GetSessionRow{}, errors.New("session not found in context")
	}
	return session, nil
}
