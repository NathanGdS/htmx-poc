package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

var tmpl *template.Template

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var tasks []Todo
var currId = 1

func init() {
	tmpl, _ = template.ParseGlob("templates/*.html")
}

func main() {
	firstTask := Todo{currId, "Teste", false}

	tasks = append(tasks, firstTask)

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("GET /tasks", fetchTodosHandler)

	http.HandleFunc("POST /tasks", addTodoHandler)

	http.HandleFunc("/task/complete", completeTodoHandler)

	http.HandleFunc("/tasks/count", countTodosHandler)

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Updated route")
		w.Header().Set("Content-Type", "text/html")
		err := tmpl.ExecuteTemplate(w, "update", nil)

		if err != nil {
			http.Error(w, "Erro ao renderizar o template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server running at http://localhost:4000")
	http.ListenAndServe(":4000", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "home.html", nil)

	if err != nil {
		http.Error(w, "Erro ao renderizar o template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func fetchTodosHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "taskList", tasks)

	if err != nil {
		http.Error(w, "Erro ao renderizar a lista de tarefas: "+err.Error(), http.StatusInternalServerError)
	}
}

func completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Complete handler")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	num, convErr := strconv.Atoi(id)

	if convErr != nil {
		log.Fatalf("Error converting string to integer: %v", convErr)
	}

	for idx, task := range tasks {
		if task.ID == num {
			tasks[idx].Completed = !tasks[idx].Completed
			w.Header().Set("HX-Trigger", "refreshCounter")
			err := tmpl.ExecuteTemplate(w, "taskItem", tasks[idx])
			log.Printf("Task ID %d updated", idx)

			if err != nil {
				http.Error(w, "Erro ao renderizar a lista de tarefas: "+err.Error(), http.StatusInternalServerError)
			}

			return
		}
	}
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	currId++
	newTask := Todo{currId, title, false}

	tasks = append(tasks, newTask)

	log.Println("new todo inserted")

	err := tmpl.ExecuteTemplate(w, "taskItem", newTask)

	if err != nil {
		http.Error(w, "Erro ao renderizar a lista de tarefas: "+err.Error(), http.StatusInternalServerError)
	}
}

func countTodosHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("counting...")
	count := len(tasks)
	var completedTasks = 0

	for _, task := range tasks {
		if task.Completed {
			completedTasks++
		}
	}

	// Executa o template criado (templates/counter.html)
	err := tmpl.ExecuteTemplate(w, "report-display", struct {
		TotalTasks     int
		CompletedTasks int
		PendingTasks   int
	}{
		TotalTasks:     count,
		CompletedTasks: completedTasks,
		PendingTasks:   count - completedTasks,
	})

	if err != nil {
		http.Error(w, "Erro ao renderizar a lista de tarefas: "+err.Error(), http.StatusInternalServerError)
	}
}
