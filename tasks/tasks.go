package tasks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"Completed"`
}

func ListTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No hay tareas")
		return
	}

	for _, task := range tasks {
		status := " "
		if task.Completed {
			status = "✓"
		}
		fmt.Printf("\n[%s] %d: %s", status, task.ID, task.Name)
	}
	fmt.Print("\n")

}

func AddTask(tasks []Task, name string) []Task {
	newTask := Task{
		ID:        generetaNextId(tasks),
		Name:      name,
		Completed: false,
	}

	return append(tasks, newTask)
}

func CompletTask(tasks []Task, id int) []Task {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = true
			break
		}
	}

	return tasks
}

func DeleteTask(tasks []Task, id int) []Task {
	for i, task := range tasks {
		if task.ID == id {
			// se contatena un slice antes del elemento en cuestión
			// con el resto de elementos del array apartir
			// de dicho elemento, pero «esparciendo» con
			// el spread operator, que aquí va detrás
			return append(tasks[:i], tasks[i+1:]...)
		}
	}

	return tasks
}

func SaveTasks(file *os.File, tasks []Task) {
	bytes, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		panic(err)
	}

	// Ponemos el puntero de escritura en el principio del fichero
	_, err = file.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	// Limpia el fichero truncando su tamaño a 0, si introducimos ese valor
	err = file.Truncate(0)
	if err != nil {
		panic(err)
	}

	// Creamos un buffer para escribir en el objeto file
	writer := bufio.NewWriter(file)
	_, err = writer.Write(bytes)
	if err != nil {
		panic(err)
	}

	// Flush se asegura de vaciar completamente el bufer de escritura
	// sobre el fichero antes de cerrarlo
	err = writer.Flush()
	if err != nil {
		panic(err)
	}

}

// Si el nombre de la función (método, struct, const, etc) comienza con
// minúsculas no será visible desde fuera del paquete; es como si fuera
// «privado»
func generetaNextId(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}

	// obtenemos el último Id y le sumamos 1
	return tasks[len(tasks)-1].ID + 1
}
