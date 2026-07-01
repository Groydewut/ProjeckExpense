package handlers

import (
	"CLIExpense/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ! Стартовая страница
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	models.ExpenseMu.Lock()
	defer models.ExpenseMu.Unlock()
	fmt.Fprint(w, "Введите запросы в строку поиска чтобы продолжить\n/expenses - посмотреть траты\n/add - добавить тарату\n/total - показать общую сумму трат")
}

// ! Созадние гет запроса, просим показать записи которые уже есть
func ExpensesHandler(w http.ResponseWriter, r *http.Request) {
	expenses, err := models.GetAllExpenses()
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expenses)
}

// !Создание пост запроса, добавление новой траты
func ExpensesCreateHandler(w http.ResponseWriter, r *http.Request) {
	//? для отправки запроса - curl -Method Post -Uri "http://localhost:8080/add" -Header @{"Content-Type"="application/json"} -Body '{"name":"Pizza","price":850}'
	var newExpense models.Expense
	err := json.NewDecoder(r.Body).Decode(&newExpense)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.ValidateExpense(newExpense)
	if response != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Ошибка при кодировании файла.", http.StatusBadRequest)
			return
		}
		return
	}

	err = models.InsertExpense(newExpense)
	if err != nil {
		var appErr models.AppError

		if errors.As(err, &appErr) {
			http.Error(w, appErr.Message, appErr.Status)
		} else {
			http.Error(w, "Не пердвиденная ошибка", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Трата добавлена")
}

// ! Создание запроса одной траты по ID
func GetExpenseByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		http.Error(w, "Отправлены не верные данные", http.StatusBadRequest)
		return
	}
	res, err := models.GetOneExpense(id)
	if err != nil {
		var appErr models.AppError

		if errors.As(err, &appErr) {
			http.Error(w, appErr.Message, appErr.Status)
		} else {
			http.Error(w, "Непредвиденная ошибка", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		http.Error(w, "Ошибка кодирования", http.StatusBadRequest)
		return
	}
}

// ! Создание DLEATE запроса
func ExpensesDel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id") // Достанет то, что попало в {id}
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		http.Error(w, "Отправлены не верные данные", http.StatusBadRequest)
		return
	}

	err = models.DeleteFromID(id)
	if err != nil {
		var appErr models.AppError

		if errors.As(err, &appErr) {
			http.Error(w, appErr.Message, appErr.Status)
		} else {
			http.Error(w, "Непрежвиденная ошибка", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Элемент удален успешно!"}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Ошибка при кодировании файла.", http.StatusBadRequest)
		return
	}
}

// ! Создание гет запроса, получение общей суммы
func TotalHandler(w http.ResponseWriter, r *http.Request) {
	total, err := models.TotalFromPrice()
	if err != nil {
		http.Error(w, "Внутриняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]float64{"total_price": total}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Ошибка при кодировании.", http.StatusInternalServerError)
		return
	}
}
