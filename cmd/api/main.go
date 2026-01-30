package main

import (
	"log"
	"net/http"
	"strings"

	"agro-monitoring/bootstrap"
	"github.com/go-chi/chi/v5"
)

func main() {
	app, err := bootstrap.NewApplication()
	if err != nil {
		log.Fatalf("Erro ao iniciar aplicação: %v", err)
	}
	defer app.Close()

	addr := ":" + app.Env.Port
	log.Printf("Servidor iniciando em http://localhost%s", addr)
	log.Printf("Endpoints disponíveis:")

	// O roteador é um http.Handler, mas para listar as rotas precisamos do tipo chi.Router
	chiRouter, ok := app.Router.(chi.Router)
	if !ok {
		log.Fatalf("O roteador não é compatível com chi.Router para listar as rotas.")
	}

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// Filtra rotas de sub-roteadores que podem conter "/*" e métodos genéricos
		if method == "*" {
			return nil
		}
		route = strings.Replace(route, "/*", "", -1)
		log.Printf("  %-6s %s", method, route)
		return nil
	}

	if err := chi.Walk(chiRouter, walkFunc); err != nil {
		log.Printf("Erro ao listar rotas: %v", err)
	}

	if err := http.ListenAndServe(addr, app.Router); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
