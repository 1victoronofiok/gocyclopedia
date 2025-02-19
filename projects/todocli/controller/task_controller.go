package controller

import (
	"errors"
	"fmt"
	"todocli/storage"
	tasks "todocli/task"
)

type TaskError struct {
	Op string
	Msg string
	Path string
}

func (te *TaskError) String() string {
	return fmt.Sprintf("error ocurred in: %s, msg: %s", te.Op, te.Msg)
}

func UpdateTaskDesc(taskId tasks.ID, desc tasks.Description) error {
	tasks, err := storage.GetAll()
	if err != nil {
		// customErr := &TaskError{Op: "Update task", Msg: "error"}
		// return fmt.Errorf("Error: %v", customErr)
		return errors.New("error: (update-task-desc, err.Error())")
	}

	tasks.UpdateTaskDesc(taskId, desc)

	storage.WriteToJsonFile(tasks)

	return nil
}

func UpdateTaskStatus(taskId tasks.ID, status tasks.Status) error {
	tasks, err := storage.GetAll()
	if err != nil {
		return fmt.Errorf("error: (update-task-status) %s", err)
	}

	err = tasks.UpdateTaskStatus(taskId, status)
	if err != nil {
		return fmt.Errorf("error: (update-task-status) %s", err)
	}

	storage.WriteToJsonFile(tasks)

	return nil
}

func DeleteTask(taskId tasks.ID) error {
	tasks, err := storage.GetAll()
	if err != nil {
		return fmt.Errorf("error: (update-task-status) %s", err)
	}

	tasks.DeleteTask(taskId)

	storage.WriteToJsonFile(tasks)

	return nil
}

func GetTasks(status tasks.Status) tasks.Tasks {
	tasks, _ := storage.GetAll()

	return tasks.GetTasksByStatus(status)
}