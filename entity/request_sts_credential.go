package entity

type RequestSTSCredential struct {
	Tenent string `json:"tenent" binding:"required"`
}
