package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func list(ctx handlerContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprintf("[Gateway][Users][List][Request] empty = %v", ""))

		listedUsers, err := ctx.UsersClient.List()
		if err != nil {
			fmt.Println(fmt.Sprintf("[Gateway][Users][List][Error] %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := Response{
			Data: listedUsers,
		}

		fmt.Println(fmt.Sprintf("[Gateway][Users][List][Response] %v", res))

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			fmt.Println(fmt.Sprintf("[Gateway][Users][List][Error] %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
