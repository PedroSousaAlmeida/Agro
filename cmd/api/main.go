package main

import (
	"log"
	"net/http"

	"agro-monitoring/bootstrap"
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
	log.Printf("  POST /v1/monitoramentos        - Upload CSV")
	log.Printf("  GET  /v1/monitoramentos        - Listar uploads")
	log.Printf("  GET  /v1/monitoramentos/:id    - Detalhes upload")
	log.Printf("  GET  /v1/areas?monitoramento_id= - Áreas de um upload")
	log.Printf("  GET  /v1/areas/:id             - Detalhes área")
	log.Printf("  GET  /v1/areas/search/fazenda  - Buscar por fazenda")
	log.Printf("  GET  /v1/areas/search/praga    - Buscar por praga")
	log.Printf("  POST /v1/areas/:id/aplicacao   - Adicionar herbicida")
	log.Printf("  GET  /health                   - Health check")

	if err := http.ListenAndServe(addr, app.Router); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
