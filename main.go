package main

import (
	"log"
	"net/http"
	"start/internal/database"
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

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err = http.ListenAndServe("localhost:7540", nil)
	if err != nil {
		panic(err)
	}
}
