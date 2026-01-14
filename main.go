package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TasksList struct {
	List []Task
}

type ListEditing interface {
	Add(Task)
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

func (t *Task) Update(s string) {
	t.Description = s
}

func (t *TasksList) Add(task Task) {
	t.List = append(t.List, task)
}

func helpText() {
	fmt.Println(`Выберете команду из списка ниже:
1. add "[description of task]"
2. update [id] "[description of task]"
2. list // all tasks
3. list done
4. list todo
5. list in-progress
`)
}

func main() {
	var tasks TasksList

	// Create JSON file with tasks
	file, err := os.Create("tasks.json")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	fmt.Println("Вас приветствует таск-менеджер.")
	helpText()

	for {
		scanCommand := bufio.NewScanner(os.Stdin)

		scanCommand.Scan()

		inputCommand := scanCommand.Text()

		splitCommand := strings.Split(inputCommand, " ")

		keyWord := splitCommand[0]

		switch keyWord {
		case "add":
			re := regexp.MustCompile(`^(\w*)\s"(.+)"$`)
			matches := re.FindStringSubmatch(inputCommand)
			if len(matches) == 0 {
				fmt.Println("Что-то пошло не так. Напишите /help для показа всех доступных команд.")
				break
			}
			tasks.Add(Task{
				ID:          uuid.New(),
				Description: string(matches[2]),
				Status:      "todo",
				CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
				UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
			})
			fmt.Println("Задача успешно добавлена!")
		case "list":
			if len(tasks.List) == 0 {
				fmt.Println("Список задач пуст!")
				helpText()
				break
			}
			fmt.Println("Список всех доступных задач: ")
			for i, v := range tasks.List {
				fmt.Printf("%d. %s \n", i+1, v.Description)
			}
		case "update":
			if len(tasks.List) == 0 {
				fmt.Println("Список задач пуст!")
				helpText()
				break
			}

			re := regexp.MustCompile(`^(\w*)\s(\d+)\s"(.+)"$`)
			matches := re.FindStringSubmatch(inputCommand)

			if len(matches) == 0 {
				fmt.Println("Что-то пошло не так. Напишите /help для показа всех доступных команд.")
				break
			}

			id, err := strconv.Atoi(matches[2])

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("ID: ", id, "LEN ", len(tasks.List))

			if id > len(tasks.List) {
				fmt.Println("Нет задачи с таким ID. Напишите /help для показа всех доступных команд.")
				break
			}

		case "/help":
			helpText()
		default:
			fmt.Println("Нет такой команды. Весь список - /help")
		}
	}

}
