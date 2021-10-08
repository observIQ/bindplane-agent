package httpevents

import "net/http"

type authMiddleware interface {
	auth(next http.Handler) http.Handler
	name() string
}

type authToken struct {
	tokenHeader string
	tokens      []string
}

func (a authToken) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(a.tokenHeader)

		for _, validToken := range a.tokens {
			if validToken == token {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusForbidden)
	})
}

func (a authToken) name() string {
	return "token-auth"
}

type authBasic struct {
	username string
	password string
}

func (a authBasic) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if ok {
			if u == a.username && p == a.password {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusForbidden)
	})
}

func (a authBasic) name() string {
	return "basic-auth"
}
