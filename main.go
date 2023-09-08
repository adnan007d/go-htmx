package main

import (
	"go-htmx/database"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

type Todo struct {
	Id        int
	Title     string
	Completed bool
}

func main() {

	database.ConnectDatabase()
	defer database.Db.Close()

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New())

	api_v1 := fiber.New()
	app.Mount("/api", api_v1)

	app.Get("/", IndexPage)

	api_v1.Post("/todos", AddTodo)
	api_v1.Post("/todos/toggle/:id", ToggleTodo)
	api_v1.Delete("/todos/:id", DeleteTodo)

	log.Fatal(app.Listen(":6969"))
}

func AddTodo(c *fiber.Ctx) error {

	title := strings.Trim(c.FormValue("title"), " ")

	if len(title) == 0 {
		c.Status(http.StatusBadRequest)
		return c.SendString("")
	}

	row := database.Db.QueryRow("INSERT INTO Todos(title, completed) values (?, ?)  RETURNING *", title, 0)

	if row.Err() != nil {
		log.Println(row.Err())
	}

	var (
		id        int
		completed bool
	)

	row.Scan(&id, &title, &completed)

	log.Println(id, title, completed)

	return c.Render("todo", fiber.Map{
		"Id":        id,
		"Title":     title,
		"Completed": completed,
	})
}

func IndexPage(c *fiber.Ctx) error {
	// order matters when using scan that is why I am specifying each column
	rows, err := database.Db.Query("SELECT id, title, completed FROM Todos")

	if err != nil {
		log.Fatal(err)
	}

	todos := []Todo{}

	for rows.Next() {
		var todo Todo
		// Fuck order matters
		rows.Scan(&todo.Id, &todo.Title, &todo.Completed)
		todos = append(todos, todo)
	}

	return c.Render("index", fiber.Map{
		"Todos": todos,
	})
}

func ToggleTodo(c *fiber.Ctx) error {

	id := c.Params("id")

	// here 1 - x it will flip from 1 and 0 hence toggle
	row := database.Db.QueryRow("UPDATE Todos SET completed = 1 - t.completed FROM (SELECT completed FROM Todos WHERE id = ?) t WHERE id = ? RETURNING *;", id, id)

	var todo Todo

	if row.Err() != nil {
		log.Println(row.Err())
	}

	row.Scan(&todo.Id, &todo.Title, &todo.Completed)

	log.Printf("%v\n", todo)

	return c.Render("todo", todo)

}

func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := database.Db.Exec("DELETE FROM Todos WHERE id = ?", id)

	if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendString("")
}
