package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	tasks "todocli/task"
)

const FILENAME string = "tasks.json"

func InitStorage(filename string) (*os.File, error) {
	_, err := os.Stat("tasks.json")

	if (os.IsExist(err)) {
		return os.Open(filename)
	}

	return nil, nil
}



/** 
	Create an instance of task with required data 
	Save (append) instance to file 
	to save, 
	retrieve data in json file, openfile with read and write flags
	convert data to Task type 
	append new data to converted data 
	convert new task data to []bytes 
	save to file 
*/
func Create(item interface {}) error {
	data, err := GetAll()
	if err != nil && err != io.EOF {
		return err
	}

	newTask, ok := item.(*tasks.Task);

	if ok  {
		fmt.Printf("NEW TASK %v, item %v ", newTask, item)
		if data == nil {
			data = &tasks.TaskStore{}
			data.Tasks = map[tasks.ID]*tasks.Task {
				0: newTask,
			}
		} else {
			data.LastId += 1
			newTask.Id = tasks.ID(data.LastId)
			data.Tasks[tasks.ID(data.LastId)] = newTask
		}
	
		WriteToJsonFile(data)
	
		fmt.Println("Written item to file")
		return nil
	}

	// fmt.Printf("\n New task type is: %v, value %v, \n", reflect.TypeOf(newTask), reflect.ValueOf(newTask))
	return fmt.Errorf("couldn't create task")
}

func DeleteById(taskId tasks.ID) error {
	file, err := openfile(FILENAME)
	if err != nil {
		return fmt.Errorf("error file opening %s: ", err)
	}
	defer file.Close()

	// unmarshal task from storage in []bytes to Task struct type and store in var (mem, should be a pointer)
	data, err := GetAll()

	if data == nil {
		return fmt.Errorf("no data")
	}

	if err != nil {
		return err
	}

	// delete task from map 
	// since it's a pointer, it should affect the initial result in mem
	if data.LastId <= 0 {
		return fmt.Errorf("error: no task entry")
	}
	
	existingTasks := data.Tasks
	delete(existingTasks, taskId)
	// data.LastId -= 1 removed this so LastId is more like the last Id

	// marshal the result to convert to []bytes (what the computer understands)
	jsonData, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		return fmt.Errorf("error marshalling %s: ", err)
	}

	err = file.Truncate(0) // Clear the file content
	if err != nil {
		return fmt.Errorf("error truncating file: %s", err)
	}
	_, err = file.Seek(0, 0) // Move the pointer to the start
	if err != nil {
		return fmt.Errorf("error seeking file: %s", err)
		
	}

	// and write back to file
    _, err = file.Write(jsonData)

	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	return nil
}

// func UpdateById(taskId tasks.ID, updatedDesc tasks.Description) error {
// 	allTask, err := GetAll()
// 	if err != nil {
// 		fmt.Printf("error %s: ", err)
// 	}

// 	// task to update holds the address of the retrieve task
// 	// so changes are persist on the original task
// 	taskToUpdate := allTask.Tasks[tasks.ID(taskId)]
// 	taskToUpdate.Description = tasks.Description(updatedDesc)

// 	fmt.Println(allTask.Tasks[tasks.ID(taskId)])

// 	WriteToJsonFile(allTask)

// 	return nil
// }

func openfile(filename string) (*os.File, error){
	return os.OpenFile(filename, os.O_RDWR | os.O_CREATE, 0644) 
}

func GetAll() (*tasks.TaskStore, error) {
	// read the file get all the data as []bytes
	file, err := openfile(FILENAME)
	if err != nil {	
		return nil, 
			fmt.Errorf("error: (GetAll-OpenFile) %s: ", err)
	}
	defer file.Close()
	
	fileBuffer, err := readFile(file)
	if err == io.EOF {
		fmt.Println("End of file error")
		return nil, err
	}

	taskStore := &tasks.TaskStore{}
	err = json.Unmarshal(fileBuffer, taskStore)

	if err != nil {	
		return nil, 
			fmt.Errorf("error: (GetAll-Unmarshall) %s: ", err)
	}

	return taskStore, nil
}


func GetTaskById(id tasks.ID) (*tasks.Task, error) {
	allTask, err := GetAll()

	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return allTask.Tasks[id], nil
}

func readFile(file *os.File) ([]byte, error) {
	buffer := make([]byte, 4092)
	n, err := file.Read(buffer)

	if err != nil {	
		return nil, err
	}

	return buffer[:n], nil
}

func WriteToJsonFile(data interface{}) error {
	file, err := openfile(FILENAME)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("error: (WriteToJson) %s", err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error: (WriteToJson) %s", err)
	}

	_, err = file.Write(jsonData)

	if err != nil {
		return fmt.Errorf("error: (WriteToJson) %s", err)
	}

	return nil
}

// func errorHandler(err error) string {
// 	return error.New(fmt.Sprintf("error: (readfile) %s: ", err))
// }

/**
	repeated tasks
	opening and closing file 
	handling or printing io errors
	unmarshalling and marshalling files
	writing to file 
	retrieving from file 
*/