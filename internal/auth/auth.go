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
const keyCacheExpiration = 1 * time.Hour
const sessionExpiration = 30 * 24 * time.Hour
const sessionExpirationSeconds = 30 * 24 * 3600
const sessionCacheExpiration = 30 * time.Minute
const cacheCleanupInterval = 15 * time.Minute
const CookieName = "sessionId"

type Service struct {
	db            *database.Queries
	sessionCache  *cache.Cache[string, database.GetSessionRow]
	apiKeyCache   *cache.Cache[[32]byte, database.GetApiKeyInfoRow]
	quoteKeyCache *cache.Cache[string, database.GetQuoteKeyInfoRow]
}

func NewAuthService(db *database.Queries) *Service {
	return &Service{
		db:            db,
		sessionCache:  cache.NewCache[string, database.GetSessionRow](sessionCacheExpiration, cacheCleanupInterval),
		apiKeyCache:   cache.NewCache[[32]byte, database.GetApiKeyInfoRow](keyCacheExpiration, cacheCleanupInterval),
		quoteKeyCache: cache.NewCache[string, database.GetQuoteKeyInfoRow](keyCacheExpiration, cacheCleanupInterval),
	}
}

func (a *Service) DeleteSession(r *http.Request) error {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return err
	}
	sessionId := cookie.Value

	a.sessionCache.Delete(sessionId)
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
	token := base64.RawStdEncoding.EncodeToString(bytes)

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
	a.sessionCache.SetWithDefaultExpiration(token, database.GetSessionRow(session))
	return token, nil
}

func (a *Service) CreateShopifyApiKey(ctx context.Context, params database.CreateShopifyApiKeyParams) ([32]byte, string, error) {
	// Create token
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return [32]byte{}, "", err
	}
	token := base64.RawStdEncoding.EncodeToString(bytes)

	// Hash the token
	var hash [32]byte = sha256.Sum256([]byte(token))
	params.ApiKey = hash[:] // Store the returned hash in params.ApiKey

	// Create quote key
	quoteBytes := make([]byte, 16)
	_, err = rand.Read(quoteBytes)
	if err != nil {
		return [32]byte{}, "", err
	}
	quoteKey := base64.RawURLEncoding.EncodeToString(quoteBytes)

	params.QuoteKey = quoteKey
	apiKeyInfo, err := a.db.CreateShopifyApiKey(ctx, params)
	if err != nil {
		return [32]byte{}, "", err
	}

	// Insert in key cache
	a.apiKeyCache.SetWithDefaultExpiration(hash, database.GetApiKeyInfoRow(apiKeyInfo))

	return hash, token, nil
}

func (a *Service) ValidateToken(ctx context.Context, token string) (database.GetSessionRow, bool) {

	if session, ok := a.sessionCache.Get(token); ok {
		return session, true
	}

	session, err := a.db.GetSession(ctx, token)
	if err != nil {
		return database.GetSessionRow{}, false
	}

	a.sessionCache.Set(token, session, max(sessionCacheExpiration, time.Until(session.ExpiresAt.Time)))
	return session, true
}

func (a *Service) ValidateApiKey(ctx context.Context, keyHash [32]byte) (database.GetApiKeyInfoRow, bool) {

	if session, ok := a.apiKeyCache.Get(keyHash); ok {
		return session, true
	}

	session, err := a.db.GetApiKeyInfo(ctx, keyHash[:])
	if err != nil {
		return database.GetApiKeyInfoRow{}, false
	}

	a.apiKeyCache.Set(keyHash, session, keyCacheExpiration)
	return session, true
}

func (a *Service) IsAuthenticated(r *http.Request) (database.GetSessionRow, bool) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return database.GetSessionRow{}, false
	}
	return a.ValidateToken(r.Context(), cookie.Value)
}
