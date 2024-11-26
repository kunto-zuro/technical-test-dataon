package service

import (
	"technical-test-dataon/dto"
	"technical-test-dataon/models"
)

type NodeService interface {
	GetTree() ([]dto.NodeResponse, error)
	GetNodeByID(nodeID uint) (dto.NodeResponse, error)
	CreateNode(node *models.Node) error
	UpdateNode(nodeID uint, updatedNode *models.Node) error
	DeleteNodeWithChildren(nodeID uint) error
	CheckDepthAndCircularReference(node *models.Node) (bool, bool, error)
	BulkInsert(nodes []dto.NodeRequest, parentID *uint) error
}
