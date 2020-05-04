package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	usersclient "github.com/microapis/users-api/client"
	usersHTTP "github.com/microapis/users-api/http"
)

var (
	httpserver *http.Server
)

func main() {
	//
	// LOAD ENVIRONMENTS
	//
	host := os.Getenv("HOST")
	if host == "" {
		log.Fatalln("missing env variable HOST")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("missing env variable PORT")
	}
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		log.Fatalln("missing env variable HTTP_PORT")
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	httpAddr := fmt.Sprintf("%s:%s", host, httpPort)

	//
	// gRPC CLIENT
	//
	uc, err := usersclient.New(addr)
	if err != nil {
		log.Fatal(err)
	}

	//
	// HTTP SERVER
	//
	r := mux.NewRouter()
	usersHTTP.Routes(r, uc)
	r.Use(loggingMiddleware)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}).Handler(r)

	log.Println("Starting HTTP service...")
	go func() {
		log.Println(fmt.Sprintf("HTTP service running, Listening on: %v", httpAddr))

		httpserver = &http.Server{Addr: httpAddr, Handler: handler}
		err := httpserver.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	//
	// LISTEN SIGNALS
	//
	quit := make(chan struct{})
	listenInterrupt(quit)
	<-quit
	gracefullShutdown()
}

func listenInterrupt(quit chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		log.Println("Signal received", s)
		quit <- struct{}{}
	}()
}

func gracefullShutdown() {
	log.Println("Gracefully shutdown")
	if err := httpserver.Shutdown(nil); err != nil {
		log.Fatalln(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}
