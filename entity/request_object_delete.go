package entity

type RequestObjectDelete struct {
	Path string `json:"path" binding:"required"`
	Type int8   `json:"type" binding:"required"` // 0: folder 1: file
}
