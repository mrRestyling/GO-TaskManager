package main

import (
	"log"
	"net/http"
	"start/internal/database"
	"start/internal/handlers"
)

// func getPort(p string) string {
// 	os.Setenv(p, )
// 	return os.Getenv("TODO_PORT")
// }

func main() {

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
			// case http.MethodGet:
			// 	handler.GetTaskByID(w, r)
		}
	})
	http.HandleFunc("/api/tasks", handler.GetTasks)

	err = http.ListenAndServe("localhost:7540", nil)
	if err != nil {
		panic(err)
	}
}
