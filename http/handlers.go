package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/microapis/users-api"
)

func listUsers(ctx *handlerContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[HTTP][Users][List][Init]")

		listedUsers, err := ctx.UsersClient.List()
		if err != nil {
			fmt.Println(fmt.Sprintf("[HTTP][Users][List][Error] %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := Response{
			Data: listedUsers,
		}

		fmt.Println(fmt.Sprintf("[HTTP][Users][List][Response] %v", res))

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			fmt.Println(fmt.Sprintf("[HTTP][Users][List][Error] %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getUser(ctx *handlerContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[Users][Get][Init]")

		userID := context.Get(r, "userID").(string)
		if userID == "" {
			err := "userID is not defined"
			fmt.Println(fmt.Sprintf("[HTTP][Get][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}

		user, err := ctx.UsersClient.Get(userID)
		if err != nil {
			fmt.Println(fmt.Sprintf("[HTTP][Users][Get][Error] %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := Response{
			Data: user,
		}

		fmt.Println(fmt.Sprintf("[HTTP][Users][Get][Response] %v", res))

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			fmt.Println(fmt.Sprintf("[HTTP][Users][Get][Error] %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func createUser(ctx *handlerContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprintf("[HTTP][Users][Create][Init]"))

		payload := &struct {
			User *users.User `json:"user"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			fmt.Println(fmt.Sprintf("[HTTP][Users][Create][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err.Error()})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}

		if payload.User == nil {
			err := "undefined user"
			fmt.Println(fmt.Sprintf("[HTTP][Users][Create][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}

		if payload.User.Email == "" || payload.User.Name == "" || payload.User.Password == "" {
			err := "undefined email, name or password"
			fmt.Println(fmt.Sprintf("[HTTP][Users][Create][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}

		fmt.Println(fmt.Sprintf("[HTTP][Users][Create][Request] payload = %v", payload))

		data, err := ctx.UsersClient.Create(payload.User)
		if err != nil {
			fmt.Println(fmt.Sprintf("[HTTP][Users][Create][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err.Error()})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}

		res := Response{
			Data: data,
		}

		fmt.Println(fmt.Sprintf("[HTTP][Users][Create][Response] %v", res))

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			fmt.Println(fmt.Sprintf("[HTTP][Users][Create][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err.Error()})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
	}
}

func updateUser(ctx *handlerContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprintf("[HTTP][Users][Update][Init]"))

		userID := context.Get(r, "userID").(string)
		if userID == "" {
			err := "userID is not defined"
			fmt.Println(fmt.Sprintf("[HTTP][Update][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}

		payload := &struct {
			User *users.User `json:"user"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			fmt.Println(fmt.Sprintf("[HTTP][Users][Update][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err.Error()})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}

		if payload.User == nil {
			err := "undefined user"
			fmt.Println(fmt.Sprintf("[HTTP][Users][Update][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}

		fmt.Println(fmt.Sprintf("[HTTP][Users][Update][Request] payload = %v", payload))

		payload.User.ID = userID

		data, err := ctx.UsersClient.Update(userID, payload.User)
		if err != nil {
			fmt.Println(fmt.Sprintf("[HTTP][Users][Update][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err.Error()})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}

		res := Response{
			Data: data,
		}

		fmt.Println(fmt.Sprintf("[HTTP][Users][Update][Response] %v", res))

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			fmt.Println(fmt.Sprintf("[HTTP][Users][Update][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err.Error()})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
	}
}

func deleteUser(ctx *handlerContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprintf("[HTTP][Users][Delete][Init]"))

		userID := context.Get(r, "userID").(string)
		if userID == "" {
			err := "userID is not defined"
			fmt.Println(fmt.Sprintf("[HTTP][Delete][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}

		err := ctx.UsersClient.Delete(userID)
		if err != nil {
			fmt.Println(fmt.Sprintf("[HTTP][Users][Delete][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err.Error()})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}

		res := Response{}

		fmt.Println(fmt.Sprintf("[HTTP][Users][Delete][Response] %v", res))

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			fmt.Println(fmt.Sprintf("[HTTP][Users][Delete][Error] %v", err))
			b, _ := json.Marshal(Response{Error: err.Error()})
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
	}
}
