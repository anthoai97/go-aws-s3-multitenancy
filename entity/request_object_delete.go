package entity

type RequestObjectDelete struct {
	Tenant string `json:"tenant" binding:"required"`
	Path   string `json:"path" binding:"required"`
}

type RequestFolderDelete struct {
	Tenant string `json:"tenant" binding:"required"`
	Path   string `json:"path" binding:"required"`
}
