package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io/fs"
	"log"
	"net/http"
	"svipp-server/assets"
	"svipp-server/internal/handlers"
	"svipp-server/internal/httputil"
	"svipp-server/internal/models"
)

func (s *server) routes() http.Handler {
	mux := chi.NewRouter()
	h := handlers.NewHandler(s.services, s.config.Domain)
	authMiddleware := NewAuthMiddleware(s.services.AuthService, s.config.Domain)

	setupBaseMiddlewares(mux)
	setupStaticServing(mux)
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		err := s.services.DBPool.Ping(r.Context())
		if err != nil {
			log.Printf("Database health check failed: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.Mount("/", setupWebRoutes(h, authMiddleware, s.config.IsProd))
	mux.Mount("/api", setupApiRoutes(h, authMiddleware, s.config.IsProd))
	mux.Mount("/api/driver", setupDriverApiRoutes(h, authMiddleware, s.config.IsProd))
	mux.Mount("/api/shopify", setupShopifyApiRoutes(h, authMiddleware, s.config.IsProd))

	// Add a catch-all route at the end
	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if httputil.IsNotJson(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		// For JSON requests, return a 404 Not Found
		httputil.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "Not Found"})
		return
	})

	return mux
}

func setupWebRoutes(h *handlers.Handler, authMiddleware *AuthMiddleware, prod bool) http.Handler {
	r := chi.NewRouter()
	r.Get("/", authMiddleware.AuthOrUnauth(h.HomePage, h.FrontPage))
	r.Get("/login", authMiddleware.RedirectIfAuthenticated(h.LoginPage))
	r.Get("/logout", h.Logout) // Add this line for the logout route
	r.Get("/orders/{uuid}", h.SingleOrderPage)

	//r.Group(func(r chi.Router) {
	//	r.Use(authMiddleware.AuthMiddleware)
	//})

	r.Group(func(r chi.Router) {
		if prod {
			r.Use(authMiddleware.AuthMiddleware, RequireRole(models.RoleAdmin))
			r.Get("/signup", h.SignupPage)
		} else {
			r.Get("/signup", authMiddleware.RedirectIfAuthenticated(h.SignupPage))
		}
	})

	return r
}

func setupShopifyApiRoutes(h *handlers.Handler, jwtMiddleware *AuthMiddleware, isProd bool) *chi.Mux {
	r := chi.NewRouter()
	r.Use(LogRequestBody)

	// Public routes (no authentication required)
	r.Post("/callback", h.ShopifyCallback)
	r.Post("/orders", h.ShopifyCallback)
	//r.Post("/webhook", h.ShopifyWebhook)

	// Add any authenticated routes here if needed
	// r.Group(func(r chi.Router) {
	//     r.Use(jwtMiddleware.JwtAuthMiddleware)
	//     // Add authenticated Shopify routes here
	// })

	return r
}

func setupApiRoutes(h *handlers.Handler, a *AuthMiddleware, isProd bool) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/auth", h.Authenticate)

	r.Group(func(r chi.Router) {
		if isProd {
			r.Use(a.AuthMiddleware, RequireRole(models.RoleAdmin))
		}
		r.Post("/users", h.CreateUser)
		r.Post("/business", h.CreateBusiness)
	})

	r.Group(func(r chi.Router) {
		r.Use(LogRequestBody)
		r.Use(a.AuthMiddleware)
		r.Get("/users/me", h.GetMyAccount)
		r.Post("/orders", h.NewOrder)
		r.Post("/orders/quote", h.GetOrderQuote)
		r.Get("/orders/my", h.GetMyOrders)
		r.Post("/orders/confirm", h.ConfirmOrder)
		r.Get("/verify-token", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	return r
}

func setupDriverApiRoutes(h *handlers.Handler, a *AuthMiddleware, isProd bool) *chi.Mux {
	r := chi.NewRouter()

	// Admin-only routes
	r.Group(func(r chi.Router) {
		if isProd {
			r.Use(a.AuthMiddleware, RequireRole(models.RoleAdmin))
		}
		r.Post("/", h.CreateDriver)
	})

	// Driver and admin routes
	r.Group(func(r chi.Router) {
		r.Use(a.AuthMiddleware, RequireRole(models.RoleDriver, models.RoleAdmin), LogRequestBody)

		r.Get("/verify-token", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Uncomment and implement these routes when ready
		// r.Get("/deliveries/my", h.GetMyDeliveries)
		// r.Post("/deliveries/{orderID}/accept", h.AcceptDelivery)
	})

	return r
}

func setupBaseMiddlewares(router *chi.Mux) {
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
}

func setupStaticServing(router *chi.Mux) {
	// Set up static file serving
	staticSubFS, err := fs.Sub(assets.EmbeddedFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	cachedFileServer := CacheControlMiddleware(http.FileServer(http.FS(staticSubFS)))
	router.Handle("/static/*", http.StripPrefix("/static/", cachedFileServer))

}
