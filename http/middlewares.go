package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// GetUserIDParam ...
func GetUserIDParam() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			vars := mux.Vars(r)
			userID := vars["user_id"]
			if userID == "" {
				err := "not found user_id pathparam"
				fmt.Println(fmt.Sprintf("[Subscriptions][Error] %v", err))
				b, _ := json.Marshal(Response{Error: err})
				http.Error(w, string(b), http.StatusForbidden)
				return
			}

			context.Set(r, "userID", userID)

			next.ServeHTTP(w, r)
		})
	}
}
