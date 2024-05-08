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

	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/task", handlers.TaskHandler)

	err = http.ListenAndServe("localhost:7540", nil)
	if err != nil {
		panic(err)
	}
}
