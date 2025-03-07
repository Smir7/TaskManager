package routes

import (
	"TaskManager/models"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/tasks", createTask)
	app.Get("/tasks", getAllTasks)
	app.Put("/tasks/:id", updateTask)
	app.Delete("/tasks/:id", deleteTask)
}

func createTask(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	task.Status = "new"
	query := "INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id"
	db := c.Locals("db").(*pgx.Conn) // Получаем подключение к БД из локальных данных

	if err := db.QueryRow(context.Background(), query, task.Title, task.Description, task.Status).Scan(&task.ID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create task"})
	}

	return c.Status(http.StatusCreated).JSON(task)
}

func getAllTasks(c *fiber.Ctx) error {
	query := "SELECT id, title, description, status, created_at, updated_at FROM tasks"
	db := c.Locals("db").(*pgx.Conn)

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch tasks"})
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not scan task"})
		}
		tasks = append(tasks, task)
	}

	return c.JSON(tasks)
}

func updateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	query := "UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4"
	db := c.Locals("db").(*pgx.Conn)

	_, err := db.Exec(context.Background(), query, task.Title, task.Description, task.Status, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update task"})
	}

	return c.Status(http.StatusOK).JSON(task)
}

func deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	query := "DELETE FROM tasks WHERE id = $1"
	db := c.Locals("db").(*pgx.Conn)

	_, err := db.Exec(context.Background(), query, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete task"})
	}

	return c.Status(http.StatusNoContent).SendString("")
}
