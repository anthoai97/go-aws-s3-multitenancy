package entity

type RequestListObjecTree struct {
	TenantID string `json:"tenant_id" binding:"required"`
	Path     string `json:"path" binding:"required"`
	Next     string `json:"filePath" binding:"required"`
}
