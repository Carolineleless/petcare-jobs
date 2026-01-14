package job

import (
	"context"
	"log"
	"time"
)

func ProcessJob(ctx context.Context, job Job) {
	log.Printf("processing job type=%s", job.Type)
	time.Sleep(2 * time.Second)
	log.Printf("job %s finished", job.ID)
}
