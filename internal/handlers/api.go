package handlers

import (
	"github.com/fizlip/microstore/internal/handlers/review_handlers"
	"github.com/fizlip/microstore/internal/handlers/vendor_handlers"
	"github.com/fizlip/microstore/internal/handlers/user_handlers"
	"github.com/fizlip/microstore/internal/handlers/product_handlers"
	"github.com/fizlip/microstore/internal/handlers/store_handlers"
	"github.com/fizlip/microstore/internal/middleware"
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handle(r *chi.Mux) {
	// Global middleware
	r.Use(chimiddle.StripSlashes)

    // PRODUCTS
	r.Route("/products", func(router chi.Router) {
		router.Get("/all", product_handlers.GetAllImages)
		router.Get("/get", product_handlers.GetProduct)
		router.Post("/search", product_handlers.Search)
        router.Get("/reviews", review_handlers.GetReviews)
	})

    // USER
	r.Route("/user", func(router chi.Router) {
		// Auth
		router.Use(middleware.Authorization)
		router.Get("/valid", user_handlers.CheckValidSession)
        // Execution 
		router.Post("/tool/run", user_handlers.RunProduct)
		// Invoice & usage
		router.Post("/invoice/upcoming", user_handlers.GetInvoice)
		router.Get("/usage", user_handlers.GetUsage)
		router.Post("/activate/billing", user_handlers.ActivateBilling)
		// Reviews
		router.Post("/review/create", review_handlers.CreateReview)
		router.Post("/review/delete", review_handlers.DeleteReview)
		router.Patch("/review/update", review_handlers.UpdateReview)
		router.Post("/review/rate", review_handlers.RateReview)
	})

    // STORE
	r.Route("/store", func(router chi.Router) {
		router.Get("/trending", store_handlers.GetTrending)
	})

    // VENDOR
	r.Route("/vendor", func(router chi.Router) {
		// Auth
		router.Use(middleware.Authorization)
		// Get set of products and their usage for current period
		router.Get("/get/usage", vendor_handlers.GetCurrentUsage)
		router.Get("/get/tools", vendor_handlers.GetTools)
        router.Post("/oci/create", vendor_handlers.CreateImage)
        router.Patch("/oci/update", vendor_handlers.UpdateImage)
        // API Keys
        router.Post("/key/create", vendor_handlers.CreateAPIKey)
        router.Delete("/key/delete", vendor_handlers.DeleteAPIKey)
        router.Get("/key/read", vendor_handlers.ReadAPIKey)
	})
}
