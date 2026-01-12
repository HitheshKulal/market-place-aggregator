package requests

type CreateMappingRequest struct {
	TemplateID uint              `json:"templateId" binding:"required"`
	FieldMap   map[string]string `json:"fieldMap" binding:"required"`
}
