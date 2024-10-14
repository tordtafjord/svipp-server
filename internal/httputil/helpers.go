package httputil

import (
	"net/http"
)

func IsNotJson(r *http.Request) bool {
	return r.Header.Get("Content-Type") != "application/json"
}
