// Package midway provides methods to arrange your HTTP middleware.
package midway

import (
	"net/http"
)

// Middleware wraps an http.Handler. Use it to insert code before or after a
// given handler calls ServeHTTP.
type Middleware func(http.Handler) http.Handler

// Queue returns a Middleware where ms are applied to h first-in-first-out. In
// this arrangement, the last middlware will execute its ServeHTTP func first.
//
// i.e. Queue(m1, m2)(handler) = m2(m1(handler)).
func Queue(ms ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		for i := range ms {
			h = ms[i](h)
		}
		return h
	}
}

// Stack returns a Middleware where ms are applied to h last-in-first-out. In
// this arrangement, the first middlware will execute its ServeHTTP func first.
//
// i.e. Stack(m1, m2)(handler) = m1(m2(handler)).
func Stack(ms ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		for i := range ms {
			h = ms[len(ms)-1-i](h)
		}
		return h
	}
}
