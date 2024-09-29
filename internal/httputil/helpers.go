package httputil

import (
	"context"
	"net/http"
)

type contextKey string

const IsJsonContextKey contextKey = "isJson"

func IsNotJson(r *http.Request) bool {
	return r.Header.Get("Content-Type") != "application/json"
}

func IsJsonFromContext(ctx context.Context) (bool, bool) {
	u, ok := ctx.Value(IsJsonContextKey).(bool)
	return u, ok
}
