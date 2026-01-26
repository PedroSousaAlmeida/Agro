package bootstrap

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	areaHandler "agro-monitoring/internal/modules/area/handler"
	jobsHandler "agro-monitoring/internal/modules/jobs/handler"
	monitoringHandler "agro-monitoring/internal/modules/monitoring/handler"
	userHandler "agro-monitoring/internal/modules/user/handler"
	sharedMiddleware "agro-monitoring/internal/shared/middleware"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(monHandler *monitoringHandler.Handler, areaHdlr *areaHandler.Handler, jobHdlr *jobsHandler.Handler, userHdlr *userHandler.UserHandler, auth *sharedMiddleware.Authenticator) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(sharedMiddleware.CORS)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API v1
	r.Route("/v1", func(r chi.Router) {
		// Public routes
		monHandler.RegisterRoutes(r)
		areaHdlr.RegisterRoutes(r)
		jobHdlr.RegisterRoutes(r)

		// Protected routes
		r.Route("/users", func(r chi.Router) {
			r.Use(auth.Auth)
			userHdlr.RegisterRoutes(r)
		})
	})

	return r
}
