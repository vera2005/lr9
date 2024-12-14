package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "catjkm8800"
	dbname   = "count"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

// Структура для валидации входящих данных
type CountInput struct {
	Val float32 `json:"val"` // Используем float32 для автоматической проверки числового значения
}

// Обработчик GET запроса
func (h *Handlers) GetCount(c echo.Context) error {
	msg, err := h.dbProvider.SelectCount()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, msg)
}

// Обработчик POST запроса
func (h *Handlers) PostCount(c echo.Context) error {
	input := CountInput{}

	// Привязка входных данных и проверка на ошибки
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "Неправильный формат JSON")
	}
	if err := h.dbProvider.InsertCount(input.Val); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, "Значение успешно вставлено")
}

// Обработчик PUT запроса
func (h *Handlers) PutCount(c echo.Context) error {
	input := CountInput{}

	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "Неправильный формат JSON")
	}

	if err := h.dbProvider.UpdateCount(input.Val); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Значение успешно обновлено")
}

// Методы для работы с базой данных
func (dbp *DatabaseProvider) SelectCount() (string, error) {
	var msg string
	row := dbp.db.QueryRow("SELECT summa FROM count ORDER BY id DESC LIMIT 1")

	err := row.Scan(&msg)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (dbp *DatabaseProvider) InsertCount(v float32) error {
	_, err := dbp.db.Exec("INSERT INTO count (val, summa) VALUES ($1, $1 + (SELECT COALESCE(summa, 0) FROM count ORDER BY id DESC LIMIT 1))", v)
	return err
}

func (dbp *DatabaseProvider) UpdateCount(v float32) error {
	_, err := dbp.db.Exec("UPDATE count SET val = $1, summa = (val + (SELECT summa FROM count WHERE id = ((SELECT MAX(id) FROM count) - 1))) WHERE id = (SELECT MAX(id) FROM count)", v)
	return err
}

func main() {
	address := flag.String("address", "127.0.0.1:8081", "адрес для запуска сервера")
	flag.Parse()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}

	fmt.Println("Подключено к базе данных!")

	dp := DatabaseProvider{db: db}
	h := Handlers{dbProvider: dp}

	e := echo.New()

	e.GET("/get", h.GetCount)
	e.POST("/post", h.PostCount)
	e.PUT("/put", h.PutCount)

	if err = e.Start(*address); err != nil {
		log.Fatal(err)
	}
}
