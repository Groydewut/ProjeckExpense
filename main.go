package main

import (
	"CLIExpense/handlers"
	"CLIExpense/models"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

// todo ВЕБ СЕРВЕР
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка чтения файла .env")
	}
	err = models.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	//!Создание главного роутера
	r := chi.NewRouter()
	//!подключаем логер middleware, после этого каждый запрос записывается в консоль
	r.Use(middleware.Logger)

	//!определяем маршруты(rest API)

	r.Get("/expenses", handlers.ExpensesHandler)
	r.Get("/expenses/{id}", handlers.GetExpenseByID)
	r.Post("/add", handlers.ExpensesCreateHandler)
	r.Get("/total", handlers.TotalHandler)
	// Магия chi: красивый URL-параметр {id} вместо ?id=...
	r.Delete("/delete/{id}", handlers.ExpensesDel)
	r.Get("/", handlers.HelloHandler)

	fmt.Println("Сревер запущен на http://localhost:8080 ")

	log.Fatal(http.ListenAndServe(":8080", r))

}
