package models

type User struct{
	Email string `json:"email"  binding:"required"`
	Name string `json:"name"  binding:"required"`
	Password string	`json:"password"  binding:"required"`
}

type Credentials struct{
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserID struct{
	Email string `json:"email" binding:"required"`
}

type UserOutput struct{
	Email string `json:"email" binding:"required"`
	Name string `json:"name" binding:"required"`
}
