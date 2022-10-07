package main

import "github.com/hibiken/asynq"

var (
	QUEUE_SERVER  *asynq.Server
	QUEUE_HANDLER *asynq.ServeMux
)

func SetupTaskQueueWorkerServer() {
	QUEUE_SERVER = asynq.NewServer(
		asynq.RedisClientOpt{Addr: ""},
		asynq.Config{
			Concurrency: 4,
			Queues: map[string]int{
				"critical": 2,
				"default":  1,
				"low":      1,
			},
		},
	)
	QUEUE_HANDLER = asynq.NewServeMux()
}

func CloseTaskQueueServer() {
	if QUEUE_SERVER != nil {
		QUEUE_SERVER.Shutdown()
	}
}
