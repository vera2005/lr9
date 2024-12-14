package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "catjkm8800"
	dbname   = "query"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

// Обработчики HTTP-запросов
func (h *Handlers) GetQuery(c echo.Context) error {
	msg, err := h.dbProvider.SelectQuery()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "Hello "+msg+"!")
}

func (h *Handlers) PostQuery(c echo.Context) error {
	nameInput := c.QueryParam("name") // Получаем Query-параметр
	if nameInput == "" {
		return c.String(http.StatusBadRequest, "Missing 'name' query parameter")
	}
	re := regexp.MustCompile(`[a-zA-Zа-яА-Я]`)
	if !re.MatchString(nameInput) {
		return c.String(http.StatusBadRequest, "empty string")
	}
	err := h.dbProvider.InsertQuery(nameInput)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusCreated, "Created")
}

func (h *Handlers) PutQuery(c echo.Context) error {
	nameInput := c.QueryParam("name") // Получаем Query-параметр
	if nameInput == "" {
		return c.String(http.StatusBadRequest, "Missing 'name' query parameter")
	}
	re := regexp.MustCompile(`[a-zA-Zа-яА-Я]`)
	if !re.MatchString(nameInput) {
		return c.String(http.StatusBadRequest, "empty string")
	}
	err := h.dbProvider.UpdateQuery(nameInput)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "Updated")
}

// Методы для работы с базой данных
func (dbp *DatabaseProvider) SelectQuery() (string, error) {
	var msg string
	row := dbp.db.QueryRow("SELECT name FROM query ORDER BY id DESC LIMIT 1")
	err := row.Scan(&msg)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (dbp *DatabaseProvider) UpdateQuery(n string) error {
	_, err := dbp.db.Exec("UPDATE query SET name = $1 WHERE id = (SELECT MAX(id) FROM query)", n)
	if err != nil {
		return err
	}
	return nil
}

func (dbp *DatabaseProvider) InsertQuery(n string) error {
	_, err := dbp.db.Exec("INSERT INTO query (name) VALUES ($1)", n)
	if err != nil {
		return err
	}
	return nil
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

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	fmt.Println("Connected!")

	dp := DatabaseProvider{db: db}
	h := Handlers{dbProvider: dp}

	e := echo.New()

	e.GET("/get", h.GetQuery)
	e.POST("/post", h.PostQuery)
	e.PUT("/put", h.PutQuery)

	err = e.Start(*address)
	if err != nil {
		log.Fatal(err)
	}
}
