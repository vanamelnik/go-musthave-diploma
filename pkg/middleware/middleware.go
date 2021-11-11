package middleware

import "net/http"

type mwFunc func(next http.Handler) http.Handler
