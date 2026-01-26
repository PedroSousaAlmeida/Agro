package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"agro-monitoring/bootstrap"
)

func main() {
	// Flags
	upCmd := flag.Bool("up", false, "Executar todas as migrations pendentes")
	downCmd := flag.Bool("down", false, "Reverter última migration")
	downAllCmd := flag.Bool("down-all", false, "Reverter todas as migrations")
	versionCmd := flag.Bool("version", false, "Mostrar versão atual das migrations")
	historyCmd := flag.Bool("history", false, "Mostrar histórico de migrations")
	stepsFlag := flag.Int("steps", 0, "Número de steps para up/down (0 = todas)")
	forceFlag := flag.Int("force", -1, "Forçar versão específica (usar com cuidado)")
	createCmd := flag.String("create", "", "Criar nova migration com nome especificado")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: migrate [opções]\n\n")
		fmt.Fprintf(os.Stderr, "Comandos:\n")
		fmt.Fprintf(os.Stderr, "  -up          Executar todas as migrations pendentes\n")
		fmt.Fprintf(os.Stderr, "  -down        Reverter última migration\n")
		fmt.Fprintf(os.Stderr, "  -down-all    Reverter todas as migrations\n")
		fmt.Fprintf(os.Stderr, "  -version     Mostrar versão atual\n")
		fmt.Fprintf(os.Stderr, "  -history     Mostrar histórico completo\n")
		fmt.Fprintf(os.Stderr, "  -steps N     Número de migrations para up/down\n")
		fmt.Fprintf(os.Stderr, "  -force V     Forçar versão V (usar com cuidado)\n")
		fmt.Fprintf(os.Stderr, "  -create NAME Criar nova migration\n")
		fmt.Fprintf(os.Stderr, "\nExemplos:\n")
		fmt.Fprintf(os.Stderr, "  migrate -up              # Executa todas pendentes\n")
		fmt.Fprintf(os.Stderr, "  migrate -up -steps 1     # Executa apenas 1 migration\n")
		fmt.Fprintf(os.Stderr, "  migrate -down            # Reverte última\n")
		fmt.Fprintf(os.Stderr, "  migrate -version         # Mostra versão atual\n")
		fmt.Fprintf(os.Stderr, "  migrate -history         # Mostra histórico\n")
		fmt.Fprintf(os.Stderr, "  migrate -create add_users # Cria nova migration\n")
	}

	flag.Parse()

	// Criar nova migration
	if *createCmd != "" {
		createMigration(*createCmd)
		return
	}

	// Carrega configurações do ambiente
	env := bootstrap.NewEnv()
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		env.DBUser, env.DBPassword, env.DBHost, env.DBPort, env.DBName, env.DBSSLMode)

	// Conexão direta para histórico
	db, err := sql.Open("postgres", env.DSN())
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}
	defer db.Close()

	// Mostrar histórico
	if *historyCmd {
		showHistory(db)
		return
	}

	// Inicializa migrate
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("Erro ao inicializar migrate: %v", err)
	}
	defer m.Close()

	// Pega versão antes da execução
	versionBefore, _, _ := m.Version()

	// Executa comando
	switch {
	case *upCmd:
		start := time.Now()
		var targetVersion uint

		if *stepsFlag > 0 {
			err = m.Steps(*stepsFlag)
		} else {
			err = m.Up()
		}

		elapsed := time.Since(start)

		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Erro ao executar migrations: %v", err)
		}

		if err == migrate.ErrNoChange {
			fmt.Println("Nenhuma migration pendente")
		} else {
			targetVersion, _, _ = m.Version()
			// Registra no histórico
			for v := versionBefore + 1; v <= targetVersion; v++ {
				name := getMigrationName(int(v))
				recordHistory(db, int(v), name, "up", int(elapsed.Milliseconds()))
			}
			fmt.Printf("Migrations executadas com sucesso! (v%d -> v%d) [%dms]\n", versionBefore, targetVersion, elapsed.Milliseconds())
		}

	case *downCmd:
		start := time.Now()
		versionBefore, _, _ := m.Version()

		if *stepsFlag > 0 {
			err = m.Steps(-*stepsFlag)
		} else {
			err = m.Steps(-1)
		}

		elapsed := time.Since(start)

		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Erro ao reverter migration: %v", err)
		}

		versionAfter, _, _ := m.Version()

		// Registra no histórico
		for v := versionBefore; v > versionAfter; v-- {
			name := getMigrationName(int(v))
			recordHistory(db, int(v), name, "down", int(elapsed.Milliseconds()))
		}

		fmt.Printf("Migration revertida! (v%d -> v%d) [%dms]\n", versionBefore, versionAfter, elapsed.Milliseconds())

	case *downAllCmd:
		start := time.Now()
		versionBefore, _, _ := m.Version()

		err = m.Down()

		elapsed := time.Since(start)

		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Erro ao reverter migrations: %v", err)
		}

		// Registra no histórico
		for v := versionBefore; v >= 1; v-- {
			name := getMigrationName(int(v))
			recordHistory(db, int(v), name, "down", int(elapsed.Milliseconds()))
		}

		fmt.Printf("Todas migrations revertidas! [%dms]\n", elapsed.Milliseconds())

	case *versionCmd:
		version, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				fmt.Println("Nenhuma migration executada ainda")
				return
			}
			log.Fatalf("Erro ao obter versão: %v", err)
		}
		fmt.Printf("Versão atual: %d\n", version)
		if dirty {
			fmt.Println("⚠️  AVISO: Estado dirty (migration falhou no meio)")
		}

	case *forceFlag >= 0:
		err = m.Force(*forceFlag)
		if err != nil {
			log.Fatalf("Erro ao forçar versão: %v", err)
		}
		fmt.Printf("Versão forçada para: %d\n", *forceFlag)

	default:
		flag.Usage()
	}
}

func createMigration(name string) {
	files, err := os.ReadDir("migrations")
	if err != nil {
		log.Fatalf("Erro ao ler diretório migrations: %v", err)
	}

	nextVersion := 1
	for _, f := range files {
		if !f.IsDir() {
			var v int
			fmt.Sscanf(f.Name(), "%d_", &v)
			if v >= nextVersion {
				nextVersion = v + 1
			}
		}
	}

	upFile := fmt.Sprintf("migrations/%03d_%s.up.sql", nextVersion, name)
	downFile := fmt.Sprintf("migrations/%03d_%s.down.sql", nextVersion, name)

	if err := os.WriteFile(upFile, []byte("-- Migration UP\n\n"), 0644); err != nil {
		log.Fatalf("Erro ao criar arquivo up: %v", err)
	}
	if err := os.WriteFile(downFile, []byte("-- Migration DOWN\n\n"), 0644); err != nil {
		log.Fatalf("Erro ao criar arquivo down: %v", err)
	}

	fmt.Printf("Migrations criadas:\n")
	fmt.Printf("  %s\n", upFile)
	fmt.Printf("  %s\n", downFile)
}

func getMigrationName(version int) string {
	files, err := filepath.Glob(fmt.Sprintf("migrations/%03d_*.up.sql", version))
	if err != nil || len(files) == 0 {
		// Tenta formato sem zero-padding
		files, err = filepath.Glob(fmt.Sprintf("migrations/%d_*.up.sql", version))
		if err != nil || len(files) == 0 {
			return fmt.Sprintf("migration_%d", version)
		}
	}

	// Extrai nome do arquivo
	base := filepath.Base(files[0])
	re := regexp.MustCompile(`^\d+_(.+)\.up\.sql$`)
	matches := re.FindStringSubmatch(base)
	if len(matches) > 1 {
		return matches[1]
	}
	return base
}

func recordHistory(db *sql.DB, version int, name, direction string, executionMs int) {
	// Verifica se tabela existe
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_name = 'migration_history'
		)
	`).Scan(&exists)

	if err != nil || !exists {
		return // Tabela ainda não existe, ignora
	}

	_, err = db.Exec(`
		INSERT INTO migration_history (version, name, direction, execution_ms)
		VALUES ($1, $2, $3, $4)
	`, version, name, direction, executionMs)

	if err != nil {
		log.Printf("Aviso: não foi possível registrar histórico: %v", err)
	}
}

func showHistory(db *sql.DB) {
	// Verifica se tabela existe
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_name = 'migration_history'
		)
	`).Scan(&exists)

	if err != nil || !exists {
		fmt.Println("Tabela de histórico não existe ainda.")
		fmt.Println("Execute: migrate -up para criar")
		return
	}

	rows, err := db.Query(`
		SELECT version, name, direction, executed_at, COALESCE(execution_ms, 0)
		FROM migration_history
		ORDER BY executed_at DESC
		LIMIT 20
	`)
	if err != nil {
		log.Fatalf("Erro ao consultar histórico: %v", err)
	}
	defer rows.Close()

	fmt.Println("\n┌─────────┬────────────────────────────────┬───────┬─────────────────────┬──────────┐")
	fmt.Println("│ Version │ Name                           │ Dir   │ Executed At         │ Time(ms) │")
	fmt.Println("├─────────┼────────────────────────────────┼───────┼─────────────────────┼──────────┤")

	count := 0
	for rows.Next() {
		var version int
		var name, direction string
		var executedAt time.Time
		var executionMs int

		rows.Scan(&version, &name, &direction, &executedAt, &executionMs)

		dirIcon := "↑ up"
		if direction == "down" {
			dirIcon = "↓ down"
		}

		// Trunca nome se muito longo
		if len(name) > 30 {
			name = name[:27] + "..."
		}

		fmt.Printf("│ %7d │ %-30s │ %-5s │ %s │ %8d │\n",
			version, name, dirIcon, executedAt.Format("2006-01-02 15:04:05"), executionMs)
		count++
	}

	if count == 0 {
		fmt.Println("│                         Nenhum registro encontrado                          │")
	}

	fmt.Println("└─────────┴────────────────────────────────┴───────┴─────────────────────┴──────────┘")

	// Mostra migrations disponíveis vs executadas
	fmt.Println("\nMigrations disponíveis:")
	showAvailableMigrations()
}

func showAvailableMigrations() {
	files, err := filepath.Glob("migrations/*.up.sql")
	if err != nil {
		return
	}

	sort.Strings(files)

	for _, f := range files {
		base := filepath.Base(f)
		parts := strings.SplitN(base, "_", 2)
		if len(parts) == 2 {
			version, _ := strconv.Atoi(parts[0])
			name := strings.TrimSuffix(parts[1], ".up.sql")
			fmt.Printf("  %03d: %s\n", version, name)
		}
	}
}
