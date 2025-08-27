package threadpool

import "log"

type Job struct {
	task func()
}

type Worker struct {
	id      int
	jobChan chan Job
}

type Pool struct {
	jobQueue chan Job
	workers  []*Worker
}

func NewPool(numWorkers int) *Pool {
	return &Pool{
		jobQueue: make(chan Job),
		workers:  make([]*Worker, numWorkers),
	}
}

func NewWorker(id int, jobChann chan Job) *Worker {
	return &Worker{
		id:      id,
		jobChan: jobChann,
	}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.jobChan {
			log.Printf("Worker %d is working on a job", w.id)
			job.task()
		}
	}()
}

func (p *Pool) AddJob(task func()) {
	p.jobQueue <- Job{task: task}
}

func (p *Pool) Start() {
	for i := 0; i < len(p.workers); i++ {
		worker := NewWorker(i, p.jobQueue)
		p.workers[i] = worker
		worker.Start()
	}
}
