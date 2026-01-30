package bootstrap

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	areaHandler "agro-monitoring/internal/modules/area/handler"
	clientsHandler "agro-monitoring/internal/modules/clients/handler"
	jobsHandler "agro-monitoring/internal/modules/jobs/handler"
	monitoringHandler "agro-monitoring/internal/modules/monitoring/handler"
	userHandler "agro-monitoring/internal/modules/user/handler"
	sharedMiddleware "agro-monitoring/internal/shared/middleware"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(
	monHandler *monitoringHandler.Handler,
	areaHdlr *areaHandler.Handler,
	jobHdlr *jobsHandler.Handler,
	userHdlr *userHandler.UserHandler,
	clientsHdlr *clientsHandler.Handler,
	auth *sharedMiddleware.Authenticator,
) http.Handler {
	r := chi.NewRouter()

	// Middlewares globais
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(sharedMiddleware.CORS)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Rota pública de registro (SEM autenticação)
	r.Post("/v1/register/{slug}", clientsHdlr.RegisterUser)

	// API v1 - rotas protegidas por autenticação
	r.Route("/v1", func(r chi.Router) {
		r.Use(auth.Auth)
		r.Use(sharedMiddleware.ExtractTenancy)

		// Rotas clients (/v1/clients/me, /v1/register/{slug})
		clientsHdlr.RegisterRoutes(r)

		// Rotas existentes (com multi-tenancy)
		// TODO: Adicionar middleware RequireClient quando estiver pronto
		monHandler.RegisterRoutes(r)
		areaHdlr.RegisterRoutes(r)
		jobHdlr.RegisterRoutes(r)
		userHdlr.RegisterRoutes(r)

		// Rotas admin (futuramente com middleware RequireAdminRole)
		r.Route("/admin", func(r chi.Router) {
			// TODO: Adicionar middleware RequireAdminRole
			clientsHdlr.RegisterAdminRoutes(r)
		})
	})

	return r
}
