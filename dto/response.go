package dto

import "time"

type NodeResponse struct {
	ID           uint           `json:"id"`
	Code         string         `json:"code"`
	Name         string         `json:"name"`
	ParentID     *uint          `json:"parent_id,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	ListDivision []NodeResponse `json:"list_division,omitempty"`
}
