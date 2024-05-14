package main

import (
	"log"
	"net/http"
	"os"
	"start/internal/database"
	"start/internal/handlers"
)

// func getPort(p string) string {
// 	os.Setenv(p, )
// 	return os.Getenv("TODO_PORT")
// }

func main() {

	// *** Задание со звездочкой, переменная окружения
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	// ***

	db, err := database.New()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	handler := &handlers.Handler{Db: db}
	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.TaskHandler(w, r)
		case http.MethodGet:
			handler.GetTaskByID(w, r)
		case http.MethodPut:
			handler.UpdateTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		}
	})
	http.HandleFunc("/api/task/done", handler.DoneTask)
	http.HandleFunc("/api/tasks", handler.GetTasks)

	err = http.ListenAndServe("localhost:"+port, nil)
	if err != nil {
		panic(err)
	}
}
