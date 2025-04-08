// Package handlers contains methods that route URLs with their
// correct function handlers 
package handlers

import (
	"github.com/fizlip/microstore/internal/handlers/reviewhandlers"
	"github.com/fizlip/microstore/internal/handlers/vendorhandlers"
	"github.com/fizlip/microstore/internal/handlers/userhandlers"
	"github.com/fizlip/microstore/internal/handlers/storehandlers"
	"github.com/fizlip/microstore/internal/handlers/f4rmerhandlers"
	"github.com/fizlip/microstore/internal/handlers/apihealth"
	"github.com/fizlip/microstore/internal/handlers/producthandlers"

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
    // Health 
    r.Route("/health", func(router chi.Router) {
        router.Get("/", apihealth.IsLive)  
    })

    // PRODUCTS
	r.Route("/products", func(router chi.Router) {
        router.Get("/reviews", reviewhandlers.GetReviews)
        router.Post("/run", f4rmerhandlers.Run)
	})

    // USER
	r.Route("/user", func(router chi.Router) {
		// Auth + auth utils
		router.Use(middleware.Authorization)
		router.Get("/valid", userhandlers.CheckValidSession)
        // Is user new?
		router.Get("/exists", userhandlers.DoesUserExist)
		router.Post("/create/empty", userhandlers.CreateEmptyRecord)
		router.Delete("/delete/analytics", userhandlers.DeleteAnalyticsRecord)
        // Execution 
		router.Post("/tool/run", userhandlers.RunProduct)
		// Invoice & usage
		router.Get("/invoice/upcoming", userhandlers.GetInvoice)
		router.Get("/usage", userhandlers.GetUsage)
		router.Get("/has/used/action", userhandlers.HasUsedAction)
		router.Post("/activate/billing", userhandlers.ActivateBilling)
		// Reviews
		router.Post("/review/create", reviewhandlers.CreateReview)
		router.Post("/review/delete", reviewhandlers.DeleteReview)
		router.Post("/review/rate", reviewhandlers.RateReview)
		router.Delete("/review/anonymize", reviewhandlers.AnonymizeReview)
		router.Delete("/review/interactions/delete", reviewhandlers.DeleteInteractions)
        // F4rmers
        router.Post("/f4rmer/create", f4rmerhandlers.Create)
        router.Get("/f4rmer/read", f4rmerhandlers.Read)
        router.Get("/f4rmer/read_specific", f4rmerhandlers.ReadSpecific)
        router.Patch("/f4rmer/update", f4rmerhandlers.Update)
		router.Delete("/f4rmer/delete", f4rmerhandlers.Delete)
		router.Post("/f4rmer/add/tool", f4rmerhandlers.AddTool)
	})

    // STORE
	r.Route("/store", func(router chi.Router) {
		router.Get("/trending", storehandlers.GetTrending)
		router.Get("/get/product", storehandlers.GetProduct)
		router.Get("/get/uti/product", storehandlers.GetUTIProduct)
		router.Get("/get/product/showcase", storehandlers.GetShowcase)
		router.Get("/get/product/summary", storehandlers.GetSummary)
		router.Get("/get/recent", storehandlers.GetRecent)
		router.Get("/get/analytics", producthandlers.GetActionAnalytics)
	})

    // VENDOR
	r.Route("/vendor", func(router chi.Router) {
		// Auth
		router.Use(middleware.Authorization)
		// Get set of products and their usage for current period
		router.Get("/get/usage", vendorhandlers.GetCurrentUsage)
		router.Get("/get/tools", vendorhandlers.GetTools)
        router.Get("/oci/get/endpoint", vendorhandlers.GetEndpoint)
        router.Post("/oci/create", vendorhandlers.CreateImage)
        router.Patch("/oci/update", vendorhandlers.UpdateImage)
        router.Post("/oci/deploy", vendorhandlers.DeployImage)
        router.Delete("/delete/deployment", vendorhandlers.DeleteDeployment)
        router.Delete("/delete/image", vendorhandlers.DeleteImage)

        router.Post("/upload/image", vendorhandlers.UploadShowcaseImage)

        // API Keys
        router.Post("/key/create", vendorhandlers.CreateAPIKey)
        router.Delete("/key/delete", vendorhandlers.DeleteAPIKey)
        router.Get("/key/read", vendorhandlers.ReadAPIKey)
	})
}
