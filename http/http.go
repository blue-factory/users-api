package http

import (
	h "net/http"

	"github.com/gorilla/mux"

	uc "github.com/microapis/users-api/client"
)

// Response ...
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Meta  interface{} `json:"meta,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

type handlerContext struct {
	UsersClient *uc.Client
}

// Routes ...
func Routes(r *mux.Router, uc *uc.Client) {
	s := r.PathPrefix("/api/v1/users").Subrouter()

	// define context
	ctx := handlerContext{
		UsersClient: uc,
	}

	// GET /api/v1/users/
	s.HandleFunc("/", list(ctx)).Methods(h.MethodGet, h.MethodOptions)
}
