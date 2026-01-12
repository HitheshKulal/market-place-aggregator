package requests

type CreateTemplateRequest struct {
	Name     string   `json:"name" binding:"required"`
	Fields   []string `json:"fields" binding:"required"`
	IsActive bool     `json:"isActive"`
}
