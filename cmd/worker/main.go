package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"agro-monitoring/bootstrap"
)

func main() {
	log.Println("Iniciando Worker...")

	app, err := bootstrap.NewApplication()
	if err != nil {
		log.Fatalf("Erro ao iniciar aplicação: %v", err)
	}
	defer app.Close()

	// Context com cancelamento para shutdown graceful
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Captura sinais de término
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Sinal de término recebido, encerrando worker...")
		cancel()
	}()

	// Inicia processamento de jobs
	log.Println("Worker rodando. Pressione Ctrl+C para encerrar.")
	app.JobsUseCase.RegisterAndProcessJobs(ctx)

	log.Println("Worker encerrado.")
}
