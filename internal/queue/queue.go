package queue

import (
	"sync"
	"time"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type BatchTask struct {
	ID        string     `json:"id"`
	URL       string     `json:"url"`
	Status    TaskStatus `json:"status"`
	Progress  int        `json:"progress"`
	Error     string     `json:"error,omitempty"`
	VideoID   string     `json:"video_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

type BatchQueue struct {
	mu      sync.Mutex
	tasks   map[string]*BatchTask
	worker  chan *BatchTask
	done    chan bool
}

var instance *BatchQueue
var once sync.Once

func GetInstance() *BatchQueue {
	once.Do(func() {
		instance = &BatchQueue{
			tasks:  make(map[string]*BatchTask),
			worker: make(chan *BatchTask, 10),
			done:   make(chan bool),
		}
		go instance.run()
	})
	return instance
}

func (q *BatchQueue) run() {
	for task := range q.worker {
		q.processTask(task)
	}
}

func (q *BatchQueue) processTask(task *BatchTask) {
	task.Status = StatusProcessing
	task.Progress = 25
	
	// Simulate processing
	time.Sleep(500 * time.Millisecond)
	task.Progress = 50
	
	time.Sleep(500 * time.Millisecond)
	task.Progress = 75
	
	// Complete
	task.Status = StatusCompleted
	task.Progress = 100
	task.VideoID = "video_" + task.ID
	task.CompletedAt = time.Now()
}

func (q *BatchQueue) AddTask(url string) *BatchTask {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	task := &BatchTask{
		ID:        generateID(),
		URL:       url,
		Status:    StatusPending,
		Progress:  0,
		CreatedAt: time.Now(),
	}
	
	q.tasks[task.ID] = task
	go func() { q.worker <- task }()
	return task
}

func (q *BatchQueue) GetTask(id string) *BatchTask {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.tasks[id]
}

func (q *BatchQueue) GetAllTasks() []*BatchTask {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	tasks := make([]*BatchTask, 0, len(q.tasks))
	for _, t := range q.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

func (q *BatchQueue) ClearCompleted() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	count := 0
	for id, t := range q.tasks {
		if t.Status == StatusCompleted || t.Status == StatusFailed {
			delete(q.tasks, id)
			count++
		}
	}
	return count
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + string(rune('a'+time.Now().UnixNano()%26))
}
