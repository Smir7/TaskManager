package main

import (
	"TaskManager/database"
	"TaskManager/routes"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	app := fiber.New()

	// Подключаемся к базе данных
	db, err := database.ConnectDB()

	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	defer db.Close(context.Background())

	// Настройка маршрутов
	routes.SetupRoutes(app)

	// Запускаем сервер
	log.Fatal(app.Listen(":3000"))
}
