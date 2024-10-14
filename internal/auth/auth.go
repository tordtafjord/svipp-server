package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"svipp-server/internal/cache"
	"svipp-server/internal/database"
	"svipp-server/internal/models"
	"time"
)

// New type for context key
type contextKey string

const SessionContextKey contextKey = "session"
const sessionExpiration = 30 * 24 * time.Hour
const sessionExpirationSeconds = 30 * 24 * 3600
const sessionCacheExpiration = 30 * time.Minute
const cacheCleanupInterval = 15 * time.Minute
const CookieName = "sessionId"

type Service struct {
	db           *database.Queries
	SessionCache *cache.Cache[string, database.GetSessionRow]
}

func NewAuthService(db *database.Queries) *Service {
	return &Service{
		db:           db,
		SessionCache: cache.NewCache[string, database.GetSessionRow](sessionCacheExpiration, cacheCleanupInterval),
	}
}

func (a *Service) DeleteSession(r *http.Request) error {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return err
	}
	sessionId := cookie.Value

	a.SessionCache.Delete(sessionId)
	err = a.db.DeleteSession(r.Context(), sessionId)
	if err != nil {
		log.Printf("Failed to delete session for token %s, %v", sessionId, err)
		return err
	}
	return nil
}

func (a *Service) CreateSession(ctx context.Context, userId int64, role models.Role) (string, error) {
	// Create token
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	token := base64.StdEncoding.EncodeToString(bytes)

	// Expires
	expiresAt := time.Now().Add(sessionExpiration)
	expiresAtTs := pgtype.Timestamptz{
		Time:  expiresAt,
		Valid: true,
	}

	session, err := a.db.InsertToken(ctx, database.InsertTokenParams{
		Token:     token,
		ExpiresAt: expiresAtTs,
		UserID:    userId,
		Role:      role.String(),
	})
	if err != nil {
		return "", err
	}
	a.SessionCache.SetWithDefaultExpiration(token, database.GetSessionRow(session))
	return token, nil
}

func (a *Service) CreateApiKey(ctx context.Context, userId int64, role models.Role) (string, error) {
	// Create token
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	token := base64.StdEncoding.EncodeToString(bytes)

	// Hash the token
	hasher := sha256.New()
	hasher.Write([]byte(token))
	apiKey := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// Set expiration to a far future date
	farFuture := time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
	expiresAtTs := pgtype.Timestamptz{
		Time:  farFuture,
		Valid: true,
	}

	session, err := a.db.InsertToken(ctx, database.InsertTokenParams{
		Token:     apiKey,
		ExpiresAt: expiresAtTs,
		UserID:    userId,
		Role:      role.String(),
	})
	if err != nil {
		return "", err
	}
	a.SessionCache.Set(token, database.GetSessionRow(session), time.Hour)
	return token, nil
}

func (a *Service) ValidateToken(ctx context.Context, token string) (database.GetSessionRow, bool) {

	if session, ok := a.SessionCache.Get(token); ok {
		return session, true
	}

	session, err := a.db.GetSession(ctx, token)
	if err != nil {
		log.Printf("Failed to fetch session from db %v", err)
		return database.GetSessionRow{}, false
	}

	a.SessionCache.Set(token, session, max(sessionCacheExpiration, time.Until(session.ExpiresAt.Time)))

	return session, true
}

func (a *Service) IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		log.Printf("User not authenticated, no cookie exists %v", err)
		return false
	}
	_, ok := a.ValidateToken(r.Context(), cookie.Value)
	return ok
}
