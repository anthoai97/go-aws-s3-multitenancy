package entity

type RequestFileUpload struct {
	Tenant   string `json:"tenant" binding:"required"`
	FilePath string `json:"filePath" binding:"required"`
}

type ResponseFileUpload struct {
	Tenant    string  `json:"tenant"`
	FilePath  string  `json:"file_path"`
	UploadUrl *string `json:"upload_url"`
}
