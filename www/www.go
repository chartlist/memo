/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-09-09
 */
package main

import (
	"log"

	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()
	app.Static("/", "./")
	app.Static("/monitor", "./monitor")
	app.Static("/charts", "./charts")
	log.Fatal(app.Listen(":3000"))
}
