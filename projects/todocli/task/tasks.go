package tasks

import (
	"fmt"
	"time"
)

type Description string
type ID int
type Status string
const (
	Done Status = "done"
	InProgress Status = "in-progress"
	Todo Status = "todo"
)

type Task struct {
	Id ID `json:"id"`
	Description Description `json:"description"`
	Status Status `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

/**
	taskStore {
		LastId -> int value
		Tasks -> 0x98010 (map addresss) -> map value
	}
*/
type Tasks map[ID]*Task

type TaskStore struct {
	LastId int `json:"lastId"`
	Tasks `json:"tasks"`
}


func NewTask(desc Description) *Task {
	return &Task{
		Description: desc,
		Status: Todo,
		CreatedAt: time.Now().UTC().Format("02 Jan 2006 15:04"),
		UpdatedAt: time.Now().UTC().Format("02 Jan 2006 15:04"),
	}
}

func (task *Task) UpdateTaskDesc(desc Description) {
	task.Description = desc
	task.UpdatedAt = time.Now().UTC().Format("02 Jan 2006 15:04")
}

// tasks is a map which is a reference type, so the receiver is already a pointer type
// meaning the method operates on instances of Tasks
func (tasks Tasks) UpdateTaskStatus(id ID, status Status) error {
	if status != Done && status != InProgress && status != Todo {
		return fmt.Errorf("error - invalid status provided")
	}
	tasks[id].Status = status
	tasks[id].UpdatedAt = time.Now().UTC().Format("02 Jan 2006 15:04")

	return nil
}

func (tasks Tasks) UpdateTaskDesc(id ID, desc Description) {
	tasks[id].Description = desc
	tasks[id].UpdatedAt = time.Now().UTC().Format("02 Jan 2006 15:04")
}

func (tasks Tasks) DeleteTask(id ID) {
	delete(tasks, id)
}

func (tasks Tasks) GetTasksByStatus(status Status) Tasks {
	filteredTasks := make(Tasks)

	if status != Done && status != InProgress && status != Todo {
		return tasks
	}

	for id, task := range tasks {
		if task.Status == status {
			filteredTasks[id] = task
		}
	}

	return filteredTasks
}