package main

import (
	"fmt"
	"net/http"

	"github.com/fizlip/microstore/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)

	var r *chi.Mux = chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"http://localhost:3000"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"PATCH", "GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{
            "Accept", 
            "Authorization", 
            "Content-Type", 
            "X-CSRF-Token", 
            "X-Username",
            "X-Provider",
            "X-Delete",
        },
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	handlers.Handle(r)

	fmt.Println("Starting GO API service on localhost:8000...")
	fmt.Println(`
      _____  _____                .__                                
    _/ ____\/  |  |_______  _____ |  |__   ____  __ __  ______ ____  
    \   __\/   |  |\_  __ \/     \|  |  \ /  _ \|  |  \/  ___// __ \ 
     |  | /    ^   /|  | \/  Y Y  \   Y  (  (_) )  |  /\___ \\  ___/ 
     |__| \____   | |__|  |__|_|  /___|  /\____/|____//____  >\___  >
               |__|             \/     \/                  \/     \/ 
  `)

	err := http.ListenAndServe("localhost:8000", r)

	if err != nil {
		log.Error(err)
	}
}
