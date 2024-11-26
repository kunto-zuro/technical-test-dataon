package service

import (
	"net/http"
	"strconv"
	"technical-test-dataon/dto"
	"technical-test-dataon/models"

	"github.com/labstack/echo/v4"
)

type NodeHandler struct {
	Service NodeService
}

func NewNodeHandler(s NodeService) *NodeHandler {
	return &NodeHandler{Service: s}
}

func (h *NodeHandler) GetTree(c echo.Context) error {
	tree, err := h.Service.GetTree()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tree)
}

func (h *NodeHandler) GetNodeByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	node, err := h.Service.GetNodeByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, node)
}

func (h *NodeHandler) CreateNode(c echo.Context) error {
	node := new(models.Node)
	if err := c.Bind(node); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	if err := h.Service.CreateNode(node); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, node)
}

func (h *NodeHandler) UpdateNode(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	node := new(models.Node)
	if err := c.Bind(node); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	if err := h.Service.UpdateNode(uint(id), node); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, node)
}

func (h *NodeHandler) DeleteNode(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	if err := h.Service.DeleteNodeWithChildren(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Node and its children deleted successfully"})
}

func (h *NodeHandler) BulkInsertHandler(c echo.Context) error {
	var nodes []dto.NodeRequest
	if err := c.Bind(&nodes); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.Service.BulkInsert(nodes, nil); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Nodes inserted successfully"})
}
