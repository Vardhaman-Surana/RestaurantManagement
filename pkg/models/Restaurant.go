package models

type Restaurant struct{
	Name string `json:"name" binding:"required"`
	Lat float64 `json:"lat" binding:"required"`
	Lng	float64 `json:"lng" binding:"required"`
	CreatorID string
	OwnerID string `json:"ownerEmailID"`
}
type RestaurantOutput struct{
	ID int `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Lat float64 `json:"lat" binding:"required"`
	Lng	float64 `json:"lng" binding:"required"`
	OwnerID string `json:"ownerEmailID" binding:"required"`
}
type ResID struct{
	ID int `json:"id"`
}
