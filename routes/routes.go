package routes

import (
	"github.com/labstack/echo/v4"
	"technical-test-dataon/service"
)

func RegisterRoutes(e *echo.Echo, handler *service.NodeHandler) {
	e.GET("/tree", handler.GetTree)
	e.GET("/tree/:id", handler.GetNodeByID)
	e.POST("/tree", handler.CreateNode)
	e.POST("/tree/bulk-insert", handler.BulkInsertHandler)
	e.PUT("/tree/:id", handler.UpdateNode)
	e.DELETE("/tree/:id", handler.DeleteNode)
}
