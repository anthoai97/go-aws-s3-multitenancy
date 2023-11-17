package entity

type RequestFileDownload struct {
	Tenant string `json:"tenant" binding:"required"`
	Path   string `json:"path" binding:"required"`
}

type ResponseFilDownload struct {
	Tenant string  `json:"tenant"`
	Path   string  `json:"path"`
	Url    *string `json:"url"`
}
