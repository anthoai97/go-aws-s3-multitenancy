package entity

type RequestFileDownload struct {
	Tenant   string `json:"tenant" binding:"required"`
	FilePath string `json:"filePath" binding:"required"`
}

type ResponseFilDownload struct {
	Tenant   string  `json:"tenant"`
	FilePath string  `json:"file_path"`
	Url      *string `json:"url"`
}
