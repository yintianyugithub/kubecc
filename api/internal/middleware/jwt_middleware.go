package middleware

import (
	"greet/api/internal/config"
	"net/http"
)

type JwtMiddleware struct {
	c *config.Config
}

func NewJwtMiddleware(c *config.Config) *JwtMiddleware {
	return &JwtMiddleware{
		c: c,
	}
}

func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		next(w, r)
	}
}
