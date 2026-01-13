package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

type TasksList struct {
	List []Task
}

type ListEditing interface {
	Add(Task) (string, error)
	// Delete(int)
}

type TaskMethods interface {
	MarkDone()
	MarkInProgress()
	Update()
}

type Task struct {
	ID          [16]byte
	Description string
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

func (t TasksList) Add(task Task) {
	pp.Println(t)
	t.List = append(t.List, task)
	fmt.Println("T AFTER LIST: ")
	pp.Println(t)
}

func main() {
	var tasks TasksList

	// Create JSON file with tasks
	file, err := os.Create("tasks.json")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	fmt.Println(`Вас приветствует таск-менеджер.
Выберете команду из списка ниже:
1. add [description of task]
2. list // all tasks
3. list done
4. list todo
5. list in-progress
`)
	scanCommand := bufio.NewScanner(os.Stdin)

	scanCommand.Scan()

	inputCommand := scanCommand.Text()

	splitCommand := strings.Split(inputCommand, " ")

	fmt.Print(splitCommand[1])

	tasks.Add(Task{
		ID:          uuid.New(),
		Description: string(splitCommand[1]),
		Status:      "todo",
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	})
	// if inputCommand == "add"
	// if inputCommand == "list" && len(sliceWithTasks) != 0 {
	// 	for i, v := range sliceWithTasks {
	// 		fmt.Print(i, v)
	// 	}
	// } else {
	// 	fmt.Print("There are no tasks")
	// }
}
