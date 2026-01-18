package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
	Delete(int)
	ReadFromJSON([]byte)
}

type TaskMethods interface {
	MarkDone()
	MarkInProgress()
	Update(string)
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
	t.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
}

func (t *Task) MarkInProgress() {
	t.Status = "in-progress"
	t.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
}

func (t *Task) MarkDone() {
	t.Status = "done"
	t.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
}

func (t *TasksList) ReadFromJSON(JSON []byte) {
	err := json.Unmarshal(JSON, t)
	if err != nil {
		fmt.Println("Ошибка при анмаршале!")
		return
	}
}

func (t *TasksList) Add(task Task) {
	t.List = append(t.List, task)
}

func (t *TasksList) Delete(id int) {
	t.List = append(t.List[:id], t.List[id+1:]...)
}

func helpText() {
	fmt.Println(`Выберете команду из списка ниже:
1. add "[description of task]"
2. update [id] "[description of task]"
2. list // show all tasks
3. list [done/todo/in-progress]
4. delete [index]
5. mark-[done/in-progress] [id]
6. exit
`)
}

func checkStatusAndPrint(tasks TasksList, status string) {
	count := 1
	for _, v := range tasks.List {
		if v.Status == status {
			fmt.Printf("%d. %s \n", count, v.Description)
			count += 1
		}
	}
	if count == 1 {
		fmt.Println("Список задач со статусом ", status, " пуст")
	}
}

func main() {
	const jsonFilePath = "./tasks.json"

	var file *os.File
	var tasks TasksList

	fmt.Println(os.Stat(jsonFilePath))

	if _, err := os.Stat(jsonFilePath); err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(jsonFilePath)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			panic(err)
		}
	} else {
		file, err = os.OpenFile(jsonFilePath, os.O_RDWR, 0)
		if err != nil {
			log.Fatal(err)
		}
		byteValue, _ := io.ReadAll(file)
		tasks.ReadFromJSON(byteValue)
	}

	defer file.Close()

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
			re := regexp.MustCompile(`^(add)\s"(.+)"$`)
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
		case "delete":
			if len(tasks.List) == 0 {
				fmt.Println("Список задач пуст!")
				helpText()
				break
			}
			re := regexp.MustCompile(`^(delete)\s(\d+)$`)
			matches := re.FindStringSubmatch(inputCommand)

			if len(matches) == 0 {
				fmt.Println("Что-то пошло не так. Напишите /help для показа всех доступных команд.")
				break
			}

			id, err := strconv.Atoi(matches[2])

			if err != nil {
				log.Fatal(err)
			}

			if len(tasks.List) < id || id-1 < 0 {
				fmt.Println("Невозможно удалить задачу под данным индексом, или же такого индекса вовсе нет!")
				break
			}

			tasks.Delete(id - 1)
			fmt.Printf("Задача под индексом %d была успешно удалена!\nВесь список команд - /help\n", id)
		case "list":
			if len(tasks.List) == 0 {
				fmt.Println("Список задач пуст!")
				helpText()
				break
			}
			re := regexp.MustCompile(`^(list)\s(.+)$`)
			matches := re.FindStringSubmatch(inputCommand)

			if len(matches) == 3 {
				switch matches[2] {
				case "done":
					checkStatusAndPrint(tasks, "done")
				case "todo":
					checkStatusAndPrint(tasks, "todo")
				case "in-progress":
					checkStatusAndPrint(tasks, "in-progress")
				default:
					fmt.Println("Нет такой команды. Весь список команд - /help")
				}
			} else if len(matches) == 0 {
				fmt.Println("Список всех задач: ")
				for i, v := range tasks.List {
					fmt.Printf("%d. %s \n", i+1, v.Description)
				}
			}
		case "mark-in-progress":
			if len(tasks.List) == 0 {
				fmt.Println("Список задач пуст!")
				helpText()
				break
			}

			re := regexp.MustCompile(`^(mark-in-progress)\s(\d+)$`)
			matches := re.FindStringSubmatch(inputCommand)

			if len(matches) == 0 {
				fmt.Println("Что-то пошло не так. Напишите /help для показа всех доступных команд.")
				break
			}

			id, err := strconv.Atoi(matches[2])

			if err != nil {
				log.Fatal(err)
			}

			if len(tasks.List) < id || id-1 < 0 {
				fmt.Println("Невозможно изменить задачу под данным индексом, или же такого индекса вовсе нет!")
				break
			}

			tasks.List[id-1].MarkInProgress()
			fmt.Printf("Задача под индексом %d была успешно изменена!\nВесь список команд - /help\n", id)
		case "mark-done":
			if len(tasks.List) == 0 {
				fmt.Println("Список задач пуст!")
				helpText()
				break
			}

			re := regexp.MustCompile(`^(mark-done)\s(\d+)$`)
			matches := re.FindStringSubmatch(inputCommand)

			if len(matches) == 0 {
				fmt.Println("Что-то пошло не так. Напишите /help для показа всех доступных команд.")
				break
			}

			id, err := strconv.Atoi(matches[2])

			if err != nil {
				log.Fatal(err)
			}

			if len(tasks.List) < id || id-1 < 0 {
				fmt.Println("Невозможно изменить задачу под данным индексом, или же такого индекса вовсе нет!")
				break
			}

			tasks.List[id-1].MarkDone()
			fmt.Printf("Задача под индексом %d была успешно изменена!\nВесь список команд - /help\n", id)
		case "update":
			if len(tasks.List) == 0 {
				fmt.Println("Список задач пуст!")
				helpText()
				break
			}

			re := regexp.MustCompile(`^(update)\s(\d+)\s"(.+)"$`)
			matches := re.FindStringSubmatch(inputCommand)

			if len(matches) == 0 {
				fmt.Println("Что-то пошло не так. Напишите /help для показа всех доступных команд.")
				break
			}

			id, err := strconv.Atoi(matches[2])

			if err != nil {
				log.Fatal(err)
			}

			if id > len(tasks.List) || id-1 < 0 {
				fmt.Println("Нет задачи с таким ID. Напишите /help для показа всех доступных команд.")
				break
			}

			textUpdated := matches[3]

			tasks.List[id-1].Update(textUpdated)
			fmt.Printf("Задача под индексом %d была успешно обновлена!\nВесь список команд - /help\n", id)
		case "/help":
			helpText()
		case "exit":
			listToBytes, err := json.Marshal(tasks)

			if err != nil {
				log.Fatal(err)
			}
			if _, err = file.Seek(0, io.SeekStart); err != nil {
				log.Fatal(err)
			}

			if err = file.Truncate(0); err != nil {
				log.Fatal(err)
			}

			file.Write(listToBytes)
			return
		default:
			fmt.Println("Нет такой команды. Весь список - /help")
		}
	}

}
