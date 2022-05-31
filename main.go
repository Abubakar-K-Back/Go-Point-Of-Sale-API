package main

import (
	"fmt"
	"os"
	db "point/of/sale/db"
	"point/of/sale/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("This is the Main File...")
	godotenv.Load(".env")
	fmt.Println(os.Getenv("MYSQL_DBNAME"))

	db.Connect(os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DBNAME"), os.Getenv("MYSQL_HOST"))

	app := fiber.New()

	//cors
	app.Use(cors.New())

	//routing
	routes.Routings(app)
	app.Listen(":3030")

}
