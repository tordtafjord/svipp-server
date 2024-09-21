package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"io/fs"
	"log"
	"net/http"
	"svipp-server/assets"
	"svipp-server/internal/handlers"
	"time"
)

func (s *server) routes() http.Handler {
	mux := chi.NewRouter()
	h := handlers.NewHandler(s.config)
	jwtMiddleware := NewJWTAuthMiddleware(s.config.JWT.SecretKey)

	setupBaseMiddlewares(mux)
	setupStaticServing(mux)
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		err := s.config.DB.DBPool.Ping(r.Context())
		if err != nil {
			log.Printf("Database health check failed: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.Mount("/", setupWebRoutes(h, jwtMiddleware, s.config.IsProd))
	mux.Mount("/api", setupApiRoutes(h, jwtMiddleware, s.config.IsProd))
	mux.Mount("/api/driver", setupDriverRoutes(h, jwtMiddleware, s.config.IsProd))
	mux.Mount("/api/shopify", setupShopifyRoutes(h, jwtMiddleware, s.config.IsProd))

	return mux
}

func setupWebRoutes(h *handlers.Handler, jwtMiddleware *JWTAuthMiddleware, prod bool) http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.HomePage)
	r.Get("/login", h.LoginPage)

	return r
}

func setupShopifyRoutes(h *handlers.Handler, jwtMiddleware *JWTAuthMiddleware, isProd bool) *chi.Mux {
	r := chi.NewRouter()

	// Public routes (no authentication required)
	r.Post("/callback", h.ShopifyCallback)
	//r.Post("/webhook", h.ShopifyWebhook)

	// Add any authenticated routes here if needed
	// r.Group(func(r chi.Router) {
	//     r.Use(jwtMiddleware.JwtAuthMiddleware)
	//     // Add authenticated Shopify routes here
	// })

	return r
}

func setupApiRoutes(h *handlers.Handler, jwtMiddleware *JWTAuthMiddleware, isProd bool) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/auth", h.Authenticate)

	r.Group(func(r chi.Router) {
		if isProd {
			r.Use(jwtMiddleware.JwtAuthMiddleware, RequireRole("admin"))
		}
		r.Post("/users", h.CreateUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware.JwtAuthMiddleware)
		r.Get("/users/me", h.GetMyAccount)
		r.Post("/orders", h.NewOrder)
		r.Get("/orders/my", h.GetMyOrders)
		r.Post("/orders/confirm", h.ConfirmOrder)
		r.Get("/verify-token", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	return r
}

func setupDriverRoutes(h *handlers.Handler, jwtMiddleware *JWTAuthMiddleware, isProd bool) *chi.Mux {
	r := chi.NewRouter()

	// Admin-only routes
	r.Group(func(r chi.Router) {
		if isProd {
			r.Use(jwtMiddleware.JwtAuthMiddleware, RequireRole("admin"))
		}
		r.Post("/", h.CreateDriver)
	})

	// Driver and admin routes
	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware.JwtAuthMiddleware, RequireRole("driver", "admin"))

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
	router.Use(httprate.LimitByRealIP(100, 1*time.Minute))
}

func setupStaticServing(router *chi.Mux) {
	// Set up static file serving
	staticSubFS, err := fs.Sub(assets.EmbeddedFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(staticSubFS))
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

}
