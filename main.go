package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	task "github.com/emarifer/go-cli-crud/tasks" // Hay que hacer un alias para que no colisione con el nombre de la variable «tasks»
	// https://scene-si.org/2018/01/25/go-tips-and-tricks-almost-everything-about-imports/
)

func main() {

	file, err := os.OpenFile("tasks.json", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	var tasks []task.Task

	info, err := file.Stat()

	if err != nil {
		panic(err)
	}

	if info.Size() != 0 {
		bytes, err := io.ReadAll(file)

		if err != nil {
			panic(err)
		}

		if err = json.Unmarshal(bytes, &tasks); err != nil {
			panic(err)
		}
	} else {
		tasks = []task.Task{}
	}

	// Impresión formateada de un struct en formato JSON:
	// https://medium.com/@adriacidre/printing-structs-in-go-d76e006404c9
	// jtasks, _ := json.MarshalIndent(tasks, "", "\t")
	// fmt.Println(string(jtasks))

	var option, idStr string

	printUsage()

	for option != "exit" {
		option, idStr = getOptions()

		switch option {
		case "list":
			task.ListTasks(tasks)
			printUsage()
		case "add":
			// Introduce en un buffer de lectura lo que el usuario
			// tipee por la entrada standard
			readerName := bufio.NewReader(os.Stdin)
			fmt.Println("Ingresa el nombre de la tarea:")
			name, _ := readerName.ReadString('\n')
			name = strings.TrimSpace(name)

			tasks = task.AddTask(tasks, name)
			task.SaveTasks(file, tasks)
			printUsage()
		case "complete":
			if idStr == "" {
				fmt.Println("Debes proporcionar un Id separado por un espacio del comando «complete» para completar")
				printUsage()
				continue
			}
			// Convertimos la entrada en un integer, porque la entrada
			// realmente es un string
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("El Id debe ser un número entero")
				printUsage()
				continue
			}

			tasks = task.CompletTask(tasks, id)
			task.SaveTasks(file, tasks)
			printUsage()
		case "delete":
			if idStr == "" {
				fmt.Println("Debes proporcionar un Id separado por un espacio del comando «delete» para eliminar")
				printUsage()
				continue
			}
			// Convertimos la entrada en un integer, porque la entrada
			// realmente es un string
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("El Id debe ser un número entero")
				printUsage()
				continue
			}

			tasks = task.DeleteTask(tasks, id)
			task.SaveTasks(file, tasks)
			printUsage()
		case "exit":
			fmt.Println("\nBye!")
			os.Exit(0)
		default:
			printUsage()
		}
	}

}

func printUsage() {
	fmt.Print("\nUso: go-cli-crud [list|add|complete|delete|exit]\n\n")
}

func getOptions() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	option, _ := reader.ReadString('\n')
	input := strings.Split(option, " ")
	if len(input) < 2 {
		return strings.TrimSpace(input[0]), ""
	}
	return strings.TrimSpace(input[0]), strings.TrimSpace(input[1])
}
