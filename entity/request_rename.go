package entity

type RequestFileRename struct {
	Tenant  string `json:"tenant" binding:"required"`
	Path    string `json:"path" binding:"required"`
	NewPath string `json:"new_path" binding:"required"`
}

type ResponseFileRename struct {
	Tenant  string `json:"tenant"`
	Path    string `json:"path"`
	NewPath string `json:"new_path" binding:"required"`
}
