package entity

type RequestSTSCredential struct {
	Tenant string `json:"tenant" binding:"required"`
}
