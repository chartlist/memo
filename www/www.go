/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-09-09
 */
package main

import (
	"log"

	"github.com/gofiber/fiber"
	"github.com/urfave/cli"
)

var WebCommand = cli.Command{
	Name:   "web",
	Action: webStart,
	Usage:  "run web start",
}

func webStart(ctx *cli.Context) {
	app := fiber.New()
	app.Static("/", "./")
	app.Static("/monitor", "./monitor")
	app.Static("/charts", "./charts")
	log.Fatal(app.Listen(":5000"))
}

func main() {
	webStart(cli.NewContext(nil, nil, nil))
}
