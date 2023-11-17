package entity

type RequestFileUpload struct {
	Tenant string `json:"tenant" binding:"required"`
	Path   string `json:"path" binding:"required"`
}

type ResponseFileUpload struct {
	Tenant         string  `json:"tenant"`
	Path           string  `json:"path"`
	UploadUrl      *string `json:"upload_url"`
	HeaderTagKey   string  `json:"header_tags_key"`
	HeaderTagValue string  `json:"header_tags_value"`
}
