package service

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"technical-test-dataon/dto"
	"technical-test-dataon/models"
)

type NodeServiceImpl struct {
	DB *gorm.DB
}

func NewNodeService(db *gorm.DB) *NodeServiceImpl {
	return &NodeServiceImpl{DB: db}
}

func (s *NodeServiceImpl) GetTree() ([]dto.NodeResponse, error) {
	var nodes []models.Node
	if err := s.DB.Find(&nodes).Error; err != nil {
		return nil, err
	}

	tree := s.buildTree(nodes, nil)
	return tree, nil
}

func (s *NodeServiceImpl) GetNodeByID(nodeID uint) (dto.NodeResponse, error) {
	var nodes []models.Node
	if err := s.DB.Find(&nodes).Error; err != nil {
		return dto.NodeResponse{}, err
	}

	var node models.Node
	if err := s.DB.First(&node, "id = ?", nodeID).Error; err != nil {
		return dto.NodeResponse{}, err
	}

	return s.buildTreeByID(nodes, nodeID), nil
}

func (s *NodeServiceImpl) CreateNode(node *models.Node) error {
	var existingNode models.Node
	if err := s.DB.First(&existingNode, "code = ?", node.Code).Error; err == nil {
		return errors.New("duplicate code is not allowed")
	}

	if node.ParentID != nil {
		depth, err := s.getDepth(*node.ParentID)
		if err != nil {
			return err
		}
		if depth >= 5 {
			return errors.New("hierarchy depth cannot exceed 5 levels")
		}
	}

	isCircular, err := s.isCircularReference(node.ParentID, node.ID)
	if err != nil {
		return err
	}
	if isCircular {
		return errors.New("circular reference is not allowed")
	}

	if err := s.DB.Create(node).Error; err != nil {
		return err
	}

	return nil
}

func (s *NodeServiceImpl) BulkInsert(nodes []dto.NodeRequest, parentID *uint) error {
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.bulkInsertWithTransaction(tx, nodes, parentID, 1); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *NodeServiceImpl) bulkInsertWithTransaction(tx *gorm.DB, nodes []dto.NodeRequest, parentID *uint, level int) error {
	if level > 5 {
		return echo.NewHTTPError(http.StatusBadRequest, "Hierarchy depth cannot exceed 5 levels")
	}

	for _, node := range nodes {
		var existingNode models.Node
		if err := tx.First(&existingNode, "code = ?", node.Code).Error; err == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Duplicate code is not allowed: "+node.Code)
		}

		newNode := models.Node{
			Code:     node.Code,
			Name:     node.Name,
			ParentID: parentID,
		}

		if err := tx.Create(&newNode).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to insert node: "+err.Error())
		}

		if len(node.ListDivision) > 0 {
			if err := s.bulkInsertWithTransaction(tx, node.ListDivision, &newNode.ID, level+1); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *NodeServiceImpl) UpdateNode(nodeID uint, updatedNode *models.Node) error {
	var node models.Node
	if err := s.DB.First(&node, "id = ?", nodeID).Error; err != nil {
		return errors.New("node not found")
	}

	if updatedNode.Code != node.Code {
		var existingNode models.Node
		if err := s.DB.First(&existingNode, "code = ?", updatedNode.Code).Error; err == nil {
			return errors.New("duplicate code is not allowed")
		}
	}

	isCircular, err := s.isCircularReference(updatedNode.ParentID, nodeID)
	if err != nil {
		return err
	}
	if isCircular {
		return errors.New("circular reference is not allowed")
	}

	node.Code = updatedNode.Code
	node.Name = updatedNode.Name
	node.ParentID = updatedNode.ParentID
	return s.DB.Save(&node).Error
}

func (s *NodeServiceImpl) DeleteNodeWithChildren(nodeID uint) error {
	return s.deleteNodeWithChildren(s.DB, nodeID)
}

func (s *NodeServiceImpl) CheckDepthAndCircularReference(node *models.Node) (bool, bool, error) {
	isCircular, err := s.isCircularReference(node.ParentID, node.ID)
	if err != nil {
		return false, false, err
	}
	return true, isCircular, nil
}

func (s *NodeServiceImpl) isCircularReference(parentID *uint, nodeID uint) (bool, error) {
	currentParentID := parentID

	for currentParentID != nil {
		if *currentParentID == nodeID {
			return true, nil
		}

		var parent models.Node
		if err := s.DB.First(&parent, "id = ?", *currentParentID).Error; err != nil {
			return false, err
		}
		currentParentID = parent.ParentID
	}

	return false, nil
}

func (s *NodeServiceImpl) buildTree(nodes []models.Node, parentID *uint) []dto.NodeResponse {
	var tree []dto.NodeResponse

	for _, node := range nodes {
		if (node.ParentID == nil && parentID == nil) || (node.ParentID != nil && parentID != nil && *node.ParentID == *parentID) {
			tree = append(tree, dto.NodeResponse{
				ID:           node.ID,
				Code:         node.Code,
				Name:         node.Name,
				ParentID:     node.ParentID,
				CreatedAt:    node.CreatedAt,
				UpdatedAt:    node.UpdatedAt,
				ListDivision: s.buildTree(nodes, &node.ID),
			})
		}
	}

	return tree
}

func (s *NodeServiceImpl) buildTreeByID(nodes []models.Node, nodeID uint) dto.NodeResponse {
	var nodeResponse dto.NodeResponse

	for _, node := range nodes {
		if node.ID == nodeID {
			nodeResponse = dto.NodeResponse{
				ID:           node.ID,
				Code:         node.Code,
				Name:         node.Name,
				ParentID:     node.ParentID,
				CreatedAt:    node.CreatedAt,
				UpdatedAt:    node.UpdatedAt,
				ListDivision: s.buildTree(nodes, &node.ID),
			}
			break
		}
	}

	return nodeResponse
}

func (s *NodeServiceImpl) deleteNodeWithChildren(db *gorm.DB, nodeID uint) error {
	var children []models.Node
	if err := db.Where("parent_id = ?", nodeID).Find(&children).Error; err != nil {
		return err
	}

	for _, child := range children {
		if err := s.deleteNodeWithChildren(db, child.ID); err != nil {
			return err
		}
	}

	if err := db.Delete(&models.Node{}, "id = ?", nodeID).Error; err != nil {
		return err
	}

	return nil
}

func (s *NodeServiceImpl) getDepth(nodeID uint) (int, error) {
	level := 0
	currentParentID := &nodeID

	for currentParentID != nil {
		var parent models.Node
		if err := s.DB.First(&parent, "id = ?", *currentParentID).Error; err != nil {
			return 0, err
		}
		level++
		currentParentID = parent.ParentID
	}

	return level, nil
}
