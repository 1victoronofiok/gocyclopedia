package main

import (
	"fmt"
	"flag"
	"strconv"
	"log"
	
	"todocli/task"
	"todocli/storage"
	"todocli/controller"
)

type Command string
const (
	AddCmd Command = "add"
	DeleteCmd Command = "delete"
	UpdateCmd Command = "update"
	UpdateDescCmd Command = "update-desc"
	UpdateStatusCmd Command = "update-status"
	ListAll Command = "list-all"
	ListAllInProgress Command = "list-all-in-progress"
	ListAllNotDone Command = "list-all-not-done"
	ListAllDone Command = "list-all-done"
)

func UseCommandCli() {

	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("no args provided")
	}

	cmd, val := flag.Arg(0), flag.Args()[1:]

	switch Command(cmd) {
		case AddCmd: 
			newTask := tasks.NewTask(tasks.Description(val[0]))

			err := storage.Create(newTask)
			if err != nil {
				fmt.Println(err)
			}
		case DeleteCmd:
			taskId, err := strconv.Atoi(val[0])
			if err != nil {
				fmt.Printf("Error %s: ", err)
			}
			controller.DeleteTask(tasks.ID(taskId))
			fmt.Println("executed delete command")
		case UpdateDescCmd:
			taskId, _ := strconv.Atoi(val[0])

			fmt.Printf("new desc %s", val[1])

			controller.UpdateTaskDesc(tasks.ID(taskId), tasks.Description(val[1]))

			fmt.Println("updated task")
		case UpdateStatusCmd:
			taskId, _ := strconv.Atoi(val[0])
			
			err := controller.UpdateTaskStatus(tasks.ID(taskId), tasks.Status(val[1]))

			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Println("updated task")
		case ListAll:
			allTask, _ := storage.GetAll()

			for taskId, task := range allTask.Tasks {
				fmt.Println(taskId, task.Description)
			}
		case ListAllInProgress:
			for _, task := range controller.GetTasks(tasks.InProgress) {
				fmt.Println(task.Id, task.Description)
			}
			fmt.Print("listing all in progress")
		case ListAllDone:
			for _, task := range controller.GetTasks(tasks.Done) {
				fmt.Println(task.Id, task.Description)
			}
			fmt.Print("listing all in done")
		case ListAllNotDone:
			for _, task := range controller.GetTasks(tasks.Todo) {
				fmt.Println(task.Id, task.Description)
			}
			fmt.Print("listing all in todo")
		default:
			fmt.Printf("default in this beetch %s", cmd)
	}
}

func main() {
	UseCommandCli()
}