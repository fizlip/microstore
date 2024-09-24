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
    router.Get("/get", GetProduct)
    router.Post("/search", Search)
  })

  r.Route("/user", func(router chi.Router){
    router.Use(middleware.Authorization)
    router.Post("/oci/create", CreateImage)     
    router.Patch("/oci/update", UpdateImage)     
    router.Post("/oci/deploy", DeployImage)     
    router.Post("/oci/taskdef/create", CreateTaskdef) 
    router.Post("/oci/cluster/create", CreateCluster) 
    router.Post("/oci/run", RunProduct)
    router.Post("/invoice/upcoming", GetInvoice)
    router.Get("/usage", GetUsage)
  })

  r.Route("/store", func(router chi.Router){
    router.Get("/trending", GetTrending)     
  })
}
