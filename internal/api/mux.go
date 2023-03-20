package api

import (
	"net/http"

	"github.com/Taratukhin/TestTaskSchoolUser/core"
	"github.com/Taratukhin/TestTaskSchoolUser/internal/api/handlers"
)

const (
	pathProfile string = "/profile"
)

func NewMux(db core.Db) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(pathProfile, checkMethodAndToken(handlers.GetUserProfile(db), db))

	return mux
}
