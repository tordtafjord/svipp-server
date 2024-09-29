package httputil

import (
	"context"
	"net/http"
	"svipp-server/internal/auth"
)

func IsNotJson(r *http.Request) bool {
	return r.Header.Get("Content-Type") != "application/json"
}

func IsJsonFromContext(ctx context.Context) (bool, bool) {
	u, ok := ctx.Value(auth.IsJsonContextKey).(bool)
	return u, ok
}
