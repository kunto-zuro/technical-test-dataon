package dto

type NodeRequest struct {
	Code         string        `json:"code"`
	Name         string        `json:"name"`
	ListDivision []NodeRequest `json:"list_division"`
}
