package entity

type RequestCreateFolder struct {
	Tenant string `json:"tenant" binding:"required"`
	Path   string `json:"path" binding:"required"`
}
