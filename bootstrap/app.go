package bootstrap

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	areaHandler "agro-monitoring/internal/modules/area/handler"
	areaRepo "agro-monitoring/internal/modules/area/repository"
	areaUsecase "agro-monitoring/internal/modules/area/usecase"
	jobsHandler "agro-monitoring/internal/modules/jobs/handler"
	jobsRepo "agro-monitoring/internal/modules/jobs/repository"
	jobsUsecase "agro-monitoring/internal/modules/jobs/usecase"
	monitoringHandler "agro-monitoring/internal/modules/monitoring/handler"
	monitoringRepo "agro-monitoring/internal/modules/monitoring/repository"
	monitoringUsecase "agro-monitoring/internal/modules/monitoring/usecase"
	"agro-monitoring/internal/services/csv"
	"agro-monitoring/internal/services/queue"
)

// Application contém todas as dependências
type Application struct {
	Env         *Env
	DB          *sql.DB
	Redis       *redis.Client
	QueueSvc    queue.Service
	Router      http.Handler
	JobsUseCase jobsUsecase.JobUseCase
}

// NewApplication cria a aplicação com todas as dependências
func NewApplication() (*Application, error) {
	env := NewEnv()

	db, err := NewDatabase(env)
	if err != nil {
		return nil, err
	}

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     env.RedisAddr(),
		Password: env.RedisPassword,
		DB:       env.RedisDB,
	})

	// Queue Service
	queueSvc := queue.NewRedisQueueService(&redis.Options{
		Addr:     env.RedisAddr(),
		Password: env.RedisPassword,
		DB:       env.RedisDB,
	})

	uuidGen := func() string {
		return uuid.New().String()
	}

	// Repositories
	monRepo := monitoringRepo.NewPostgresRepository(db)
	areaRepository := areaRepo.NewPostgresRepository(db)
	jobRepository := jobsRepo.NewPostgresRepository(db)

	// Parser
	csvParser := csv.NewParser(uuidGen)

	// Use cases
	monUC := monitoringUsecase.NewMonitoringUseCase(monRepo, areaRepository, csvParser, uuidGen)
	areaUC := areaUsecase.NewAreaQueryUseCase(areaRepository)
	jobUC := jobsUsecase.NewJobUseCase(jobsUsecase.Config{
		UUIDGenerator: uuidGen,
		JobRepo:       jobRepository,
		AreaRepo:      areaRepository,
		Queue:         queueSvc,
	})

	// Handlers
	monHandler := monitoringHandler.NewHandler(monUC)
	areaHdlr := areaHandler.NewHandler(areaUC)
	jobHdlr := jobsHandler.NewHandler(jobUC)

	// Router
	router := SetupRoutes(monHandler, areaHdlr, jobHdlr)

	return &Application{
		Env:         env,
		DB:          db,
		Redis:       redisClient,
		QueueSvc:    queueSvc,
		Router:      router,
		JobsUseCase: jobUC,
	}, nil
}

// Close fecha recursos
func (app *Application) Close() error {
	if app.QueueSvc != nil {
		app.QueueSvc.Close()
	}
	if app.Redis != nil {
		app.Redis.Close()
	}
	if app.DB != nil {
		return app.DB.Close()
	}
	return nil
}
