package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Taratukhin/TestTaskSchoolUser/core"
)

const parameterName = "username"

func GetUserProfile(
	db core.Db,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			_ = r.Body.Close()
		}()

		var (
			username      string
			onlyOneObject bool
			result        []byte
		)

		if keys, ok := r.URL.Query()[parameterName]; ok { // we accept only one parameter of username
			username = keys[0]
			onlyOneObject = true
		}

		info, err := db.SelectAboutUser(r.Context(), username)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"error":"%v"}`, err)

			return
		}

		if onlyOneObject {
			if len(info) < 1 {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = fmt.Fprintf(w, `{"error":"username %s not found"}`, username)

				return
			}

			result, err = json.Marshal(info[0]) // we return only one object, not a slice
		} else {
			result, err = json.Marshal(info) // we return the slice, even if it is empty
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"error":"error to encode, %v"}`, err)

			return
		}

		_, _ = w.Write(result)
	}
}
