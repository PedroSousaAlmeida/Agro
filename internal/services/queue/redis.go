package queue

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"agro-monitoring/internal/modules/jobs/domain"
)

type RedisQueueService struct {
	client *redis.Client
}

// NewRedisQueueService cria um novo serviço de fila com Redis
func NewRedisQueueService(opts *redis.Options) Service {
	rdb := redis.NewClient(opts)
	return &RedisQueueService{client: rdb}
}

// Enqueue adiciona um job à fila
func (s *RedisQueueService) Enqueue(ctx context.Context, job *Job, opts *EnqueueOptions) error {
	queueName := QueueDefault
	if opts != nil && opts.QueueName != "" {
		queueName = opts.QueueName
	}

	payload, err := json.Marshal(job.JobEntity)
	if err != nil {
		return err
	}

	return s.client.LPush(ctx, queueName, payload).Err()
}

// Dequeue busca um job da fila de forma bloqueante
func (s *RedisQueueService) Dequeue(ctx context.Context, queueName string) (*Job, error) {
	// BRPOP bloqueia até que um item esteja disponível ou o timeout ocorra
	// Timeout de 1 segundo para permitir verificação do context cancelado
	result, err := s.client.BRPop(ctx, time.Second, queueName).Result()
	if err != nil {
		// Timeout sem job na fila - retorna nil sem erro
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	// BRPop retorna um array com [nome_da_fila, valor]
	if len(result) != 2 {
		return nil, redis.Nil
	}

	var jobEntity domain.Job
	if err := json.Unmarshal([]byte(result[1]), &jobEntity); err != nil {
		return nil, err
	}

	job := &Job{
		ID:        jobEntity.ID,
		Queue:     result[0], // O nome da fila da qual o job foi retirado
		Payload:   []byte(result[1]),
		JobEntity: &jobEntity,
	}

	return job, nil
}

// Close fecha a conexão com o Redis
func (s *RedisQueueService) Close() error {
	return s.client.Close()
}
