package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/shaural/GlobePersonalWebsite/server/pkg/api"
	"github.com/shaural/GlobePersonalWebsite/server/pkg/common"
)

func main() {
	config := common.Config()
	port := config.Port
	log.Printf("Listening on port:%s", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)
	defer close(sigChan)

	router := mux.NewRouter()

	api.AddMainHandler(router)
	api.AddMapHandler(router)

	mux := http.NewServeMux()
	mux.Handle("/", router)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sig := <-sigChan
		log.Printf("Received signal %s", sig)
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutrdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	http.Handle("/", router)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Listen and serve error: %v", err)
	}

	<-idleConnsClosed
}
