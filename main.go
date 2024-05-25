package main

import (
	"log"
	"net/http"
	"os"
	"start/internal/handlers"
	"start/internal/middleware"
	"start/internal/storage"
)

func main() {

	// * Переменная окружения
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// * Аутентификация по паролю
	password := os.Getenv("TODO_PASSWORD")

	// Создаем подключение к БД
	db, err := storage.New()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Cоздаем новый экземпляр обработчика handler и передаем ему ссылку на подключение к БД
	// файл addTask.go, каталог handlers
	handler := &handlers.Handler{
		Db:       db,
		Password: password,
	}

	if password != "" {
		log.Printf("Password: %s", password)
		http.HandleFunc("/api/signin", handler.LoginSign)
	}

	// Регистрируем обработчик маршрутов (связываем фронт и бэк)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// Регистрируем обработчик для получения следующей даты для проверки функции NextDate
	// файл getHandler.go находится в каталоге handlers
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)

	// 1.Регистрируем обработчик для обработки маршрута /api/task
	// 2.Оборачиваем обработчики в функцию middleware.Auth для авторизации по токену
	// middleware.Auth - файл mid.go, каталог middleware
	http.HandleFunc("/api/task", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.TaskHandler(w, r) // файл addTask.go, каталог handlers
		case http.MethodGet:
			handler.GetTaskByID(w, r) // файл getTaskByID.go, каталог handlers
		case http.MethodPut:
			handler.UpdateTask(w, r) // файл updateTask.go, каталог handlers
		case http.MethodDelete:
			handler.DeleteTask(w, r) // файл deleteTask.go, каталог handlers
		}
	}))

	// Обработчик для POST-запроса, который делает задачу выполненной
	http.HandleFunc("/api/task/done", middleware.Auth(handler.DoneTask)) // файл done.go, каталог handlers

	// Получаем список ближайших задач
	http.HandleFunc("/api/tasks", middleware.Auth(handler.GetTasks)) // файл getAllTasks.go, каталог handlers

	// Запускаем сервер на указанном порту
	err = http.ListenAndServe("0.0.0.0:"+port, nil)
	// err = http.ListenAndServe("localhost:"+port, nil)
	if err != nil {
		panic(err)
	}
}
