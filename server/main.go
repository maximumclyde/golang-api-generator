package main

import (
	db "api-generator/db"
	"api-generator/router"
	"api-generator/store"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//#region main
	database, err := db.OpenDb(nil)
	if err != nil {
		panic(err)
	}

	store := store.NewStore(database)
	router := router.NewRouter(store)
	srv := &http.Server{
		Addr:    "",
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down...")

	db.Close(database)

	srv.Shutdown(context.Background())
}
