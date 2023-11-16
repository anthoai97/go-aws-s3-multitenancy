package entity

type RequestFileDownload struct {
	FileName string `json:"fileName" binding:"required"`
	FilePath string `json:"filePath" binding:"required"`
}
