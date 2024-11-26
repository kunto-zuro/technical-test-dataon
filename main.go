package main

import (
	"github.com/labstack/echo/v4"
	"technical-test-dataon/config"
	"technical-test-dataon/routes"
	"technical-test-dataon/service"
)

func main() {
	db := config.InitDB()

	nodeService := service.NewNodeService(db)
	nodeHandler := service.NewNodeHandler(nodeService)

	e := echo.New()

	routes.RegisterRoutes(e, nodeHandler)

	e.Logger.Fatal(e.Start(":8123"))
}
