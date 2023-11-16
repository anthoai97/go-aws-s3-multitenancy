package entity

type RequestFileUpload struct {
	FileName string `json:"fileName" binding:"required"`
	FilePath string `json:"filePath" binding:"required"`
}

type ResponseFileUpload struct {
	FileName  string  `json:"file_name"`
	FilePath  string  `json:"file_path"`
	UploadUrl *string `json:"upload_url"`
}
