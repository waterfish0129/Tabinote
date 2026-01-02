package dto

type UserLoginDTO struct {
	Name     string `json:"name" binding:"required,email" message:"必填哦" email_err:"格式要信箱"`
	Password string `json:"password" binding:"required" message:"必填哦"`
}
