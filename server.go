package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "User"
	}
	log.Printf("Received request for %s\n", name)
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
}

func main() {
	// gorilla mux for start server
	r := mux.NewRouter()
	// processing
	r.HandleFunc("/", handler)

	// server
	srv := &http.Server{
		Handler:      r,
		Addr:         ":5000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// work file_path
	LOG_FILE_PATH := os.Getenv("LOG_FILE_PATH")
	if LOG_FILE_PATH != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   LOG_FILE_PATH,
			MaxSize:    200,
			MaxBackups: 3,
			Compress:   true,
		})
	}
	// start server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	// execution service
	gracefullShutdown(srv)
}

func gracefullShutdown(srv *http.Server) {
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interruptCh

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Shutting down")
	os.Exit(0)
}
