package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/Carolineleless/petcare-jobs/internal/job"
	"github.com/Carolineleless/petcare-jobs/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	cfg := repository.LoadPostgresConfig()
	log.Println(cfg)
	db, err := repository.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	queue := job.NewQueue(100)

	queue.Enqueue(job.Job{
		ID:   "1",
		Type: "GENERATE_PET_REPORT",
	})

	workerPool := job.NewWorkerPool(queue, 5)

	workerPool.Start(ctx)

	log.Println("PetCare Jobs API started")

	<-ctx.Done()
	log.Println("shutting down")
}
