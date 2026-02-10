package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/Carolineleless/petcare-jobs/internal/animal"
	"github.com/Carolineleless/petcare-jobs/internal/config"
	"github.com/Carolineleless/petcare-jobs/internal/http"
	"github.com/Carolineleless/petcare-jobs/internal/job"
	"github.com/Carolineleless/petcare-jobs/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	cfg := config.Load()

	db, err := repository.NewGormConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&animal.Animal{}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	animalRepo := animal.NewRepository(db)

	animalService := animal.NewAnimalService(animalRepo)

	animalHandler := animal.NewHandler(animalService)

	router := http.NewRouter(animalHandler)

	server := http.NewServer(router, cfg.Server.Port)

	queue := job.NewQueue(cfg.Job.QueueBuffer)
	workerPool := job.NewWorkerPool(queue, cfg.Job.WorkerPoolSize)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	workerPool.Start(ctx)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Println("PetCare Jobs API started")

	<-ctx.Done()

	log.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Shutdown complete")
}
