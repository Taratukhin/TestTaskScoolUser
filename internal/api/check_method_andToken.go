package api

import (
	"net/http"

	"github.com/Taratukhin/TestTaskSchoolUser/core"
	"golang.org/x/exp/slices"
)

func checkMethodAndToken(
	next http.HandlerFunc,
	db core.Db,
	allowedHTTPMethods ...string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if len(allowedHTTPMethods) == 0 {
			allowedHTTPMethods = []string{
				http.MethodGet,
			}
		}

		if !slices.Contains(allowedHTTPMethods, r.Method) {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		token := r.Header.Get("Api-key")
		if len(token) == 0 {
			w.WriteHeader(http.StatusForbidden) // although a status 401 (http.StatusUnauthorized) would be more suitable

			return
		}

		_, err := db.SelectIDByAPIKey(r.Context(), token)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)

			return
		}

		next(w, r)
	}
}
