package handlers

import (
  "github.com/go-chi/chi"
  chimiddle "github.com/go-chi/chi/middleware"
  "github.com/fizlip/microstore/internal/middleware"
)

func Handle(r *chi.Mux) {
  // Global middleware
  r.Use(chimiddle.StripSlashes)
  
  r.Route("/account", func(router chi.Router) {
    // Middleware for /account route
    router.Use(middleware.Authorization)

    router.Get("/services", GetUserServices)

  })

  r.Route("/oci", func(router chi.Router) {
    router.Get("/all", GetAllImages)
  })

  r.Route("/user", func(router chi.Router){
    router.Use(middleware.Authorization)
    router.Post("/create/oci", CreateImage)     
    router.Patch("/update/oci", UpdateImage)     
  })
}
