package main

import (
	"net/http"
	"svipp-server/cmd/api/handlers"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func (s *server) routes() http.Handler {
	h := handlers.NewHandler(s.config)
	jwtMiddleware := NewJWTAuthMiddleware(s.config.JWT.SecretKey)

	mux := chi.NewRouter()
	setupBaseMiddlewares(mux)

	setupStaticServing(mux)

	// apiRouter
	apiRouter := chi.NewRouter()
	setupApiRoutes(apiRouter, h, jwtMiddleware, s.config.IsProd)

	// Driver-specific router
	driverRouter := chi.NewRouter()
	//setupDriverRoutes(driverRouter, handler)

	// Mounting routers
	mux.Mount("/api", apiRouter)
	mux.Mount("/api/driver", driverRouter)

	return mux
}

func setupApiRoutes(router *chi.Mux, h *handlers.Handler, jwtMiddleWare *JWTAuthMiddleware, isProd bool) {
	router.Post("/auth", h.Authenticate)

	// Conditionally add the /users endpoint
	router.Group(func(r chi.Router) {
		if isProd {
			r.Use(jwtMiddleWare.JwtAuthMiddleware)
			r.Use(RequireRole("admin"))
		}
		router.Post("/users", h.CreateUser)
	})

	// Authenticated Group
	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleWare.JwtAuthMiddleware)
		r.Get("/users/me", h.GetMyAccount)
		r.Post("/orders", h.NewOrder)
		r.Get("/orders/my", h.GetMyOrders)
		r.Post("/orders/confirm", h.ConfirmOrder)
		r.Get("/verify-token", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})
}

func setupBaseMiddlewares(router *chi.Mux) {
	router.Use(middleware.Logger)
	router.Use(httprate.LimitByIP(100, 1*time.Minute))
	router.Use(middleware.Heartbeat("/health"))
	router.Use(middleware.Recoverer)
}

func setupStaticServing(router *chi.Mux) {
	// Set up static file serving
	//fileServer := http.FileServer(http.Dir("./static"))
	//router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Serve images publicly
	//imageServer := http.FileServer(http.Dir("./assets/static/images"))
	//router.Handle("/images/*", http.StripPrefix("/images/", imageServer))

	// Serve index.html for the root path
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/static/index.html")
	})

}
