// Package handlers contains methods that route URLs with their
// correct function handlers 
package handlers

import (
	"github.com/fizlip/microstore/internal/handlers/review_handlers"
	"github.com/fizlip/microstore/internal/handlers/vendor_handlers"
	"github.com/fizlip/microstore/internal/handlers/user_handlers"
	"github.com/fizlip/microstore/internal/handlers/product_handlers"
	"github.com/fizlip/microstore/internal/handlers/store_handlers"
	"github.com/fizlip/microstore/internal/handlers/f4rmer_handlers"

	"github.com/fizlip/microstore/internal/middleware"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

// Handle is the main HTTP router of the f4rmhouse API. 
// It uses the [chi] router to map incoming requests with their 
// correct http handler
//
// we support and will continue to add many handlers. these handlers 
// are divided into sub-packages named after the entity 
// that they handle, e.g. user_handler is the 
// package that handles all user endpoints etc. 
//
// Some endpoints (the majority actually) require authentication
// this is defined in the middleware package. (you find the all middleware in internal/middleware)
// 
// [chi]: https://github.com/go-chi/chi
func Handle(r *chi.Mux) {
	// Global middleware
	r.Use(chimiddle.StripSlashes) // this will remove unnecessary characters from the URL

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
        // F4rmers
        router.Post("/f4rmer/create", f4rmer_handlers.Create)
        router.Get("/f4rmer/read", f4rmer_handlers.Read)
        router.Get("/f4rmer/read_specific", f4rmer_handlers.ReadSpecific)
        router.Patch("/f4rmer/update", f4rmer_handlers.Update)
		router.Delete("/f4rmer/delete", f4rmer_handlers.Delete)
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
        router.Post("/oci/deploy", vendor_handlers.DeployImage)
        router.Get("/oci/get/endpoint", vendor_handlers.GetEndpoint)
        // API Keys
        router.Post("/key/create", vendor_handlers.CreateAPIKey)
        router.Delete("/key/delete", vendor_handlers.DeleteAPIKey)
        router.Get("/key/read", vendor_handlers.ReadAPIKey)
	})
}
