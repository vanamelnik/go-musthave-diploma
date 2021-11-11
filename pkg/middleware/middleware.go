package middleware

import "net/http"

type MwFunc func(next http.Handler) http.Handler
