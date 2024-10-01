package handlers

import (
  "github.com/go-chi/chi"
  chimiddle "github.com/go-chi/chi/middleware"
  "github.com/fizlip/microstore/internal/middleware"
  "github.com/fizlip/microstore/internal/handlers/review_handlers"
  "github.com/fizlip/microstore/internal/handlers/vendor_handlers"
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
    // Auth
    router.Use(middleware.Authorization)
    // OCI creation & reading & updating
    router.Post("/oci/create", CreateImage)     
    router.Patch("/oci/update", UpdateImage)     
    router.Post("/oci/deploy", DeployImage)     
    router.Post("/oci/taskdef/create", CreateTaskdef) 
    router.Post("/oci/cluster/create", CreateCluster) 
    router.Post("/oci/run", RunProduct)
    // Invoice & usage
    router.Post("/invoice/upcoming", GetInvoice)
    router.Get("/usage", GetUsage)
    // Reviews
    router.Post("/review/create", review_handlers.CreateReview)
    router.Post("/review/delete", review_handlers.DeleteReview)
    router.Post("/review/update", review_handlers.UpdateReview)
    router.Post("/review/rate", review_handlers.RateReview)
  })

  r.Route("/store", func(router chi.Router){
    router.Get("/trending", GetTrending)     
    router.Get("/product/reviews", review_handlers.GetReviews)     
  })

  r.Route("/vendor", func(router chi.Router){
    // Auth
    router.Use(middleware.Authorization)
    // Get set of products and their usage for current period
    router.Get("/get/usage", vendor_handlers.GetCurrentUsage)     
  })
}
