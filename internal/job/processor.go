package job

type Queue struct {
	jobs chan Job
}

func NewQueue(buffer int) *Queue {
	return &Queue{
		jobs: make(chan Job, buffer),
	}
}

func (q *Queue) Enqueue(job Job) {
	q.jobs <- job
}

func (q *Queue) Dequeue() <-chan Job {
	return q.jobs
}
