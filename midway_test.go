package midway

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const errorString = "\nGot:\t%v\nWant:\t%v\n"

func TestQueueAndStack(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("core"))
	})
	// sign is a Middleware that writes msg into w. In this testing context, it
	// is used to reveal the sequence of middleware execution.
	sign := func(msg string) Middleware {
		return func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(msg))
				h.ServeHTTP(w, r)
			})
		}
	}
	ms := []Middleware{sign("m1_"), sign("m2_")}
	cases := map[string]struct {
		handler http.Handler
		want    string
	}{
		"queue": {
			handler: Queue(ms...)(base),
			want:    "m2_m1_core",
		},
		"stack": {
			handler: Stack(ms...)(base),
			want:    "m1_m2_core",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
			if got := w.Body.String(); got != c.want {
				t.Errorf(errorString, got, c.want)
			}
		})
	}
}
