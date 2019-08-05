package database

import (
	"errors"
	"github.com/vds/RestaurantManagement/pkg/models"
)

//errors
var ErrInternal = errors.New("internal server error")
var ErrDupEmail=errors.New("email already used try a different one")
var ErrInvalidCredentials = errors.New("incorrect login details")


type Database interface{
	CreateUser(userType string,user *models.User)error
	LogInUser(userType string,cred *models.Credentials)(string,error)

	CreateOwners(creatorID string,owners []models.User)error
	RemoveOwners(creatorID string,ownerIDs []models.UserID)error
	RemoveOwnersBySuper(ownerIDs []models.UserID)error
	CheckOwnerCreator(creatorID string,ownerIDs []models.UserID)error
	ShowOwners(creatorID string)(string,error)
	ShowOwnersForSuper()(string,error)

	ShowRestaurants(creatorID string)(string,error)
	ShowRestaurantsForSuper()(string,error)
	InsertRestaurant(restaurant *models.Restaurant)error
	RemoveRestaurantsBySuper(resIDs []models.ResID)error
	RemoveRestaurants(creatorID string,resIDs []models.ResID)error
	CheckRestaurantCreator(creatorID string,resIDs []models.ResID)error
	UpdateRestaurant( *models.RestaurantOutput)error


	InsertDishes(dishes []models.Dish,resID int)error
	RemoveDishes(ids []models.DishID)error
	ShowMenu(resID int)(string,error)
	UpdateDishes(dishes []models.DishOutput)error

	GetOwnerRestaurants(ownerID string)(string,error)
	CheckRestaurantOwner(ownerID string,resID int)error
	IsExistingOwner(ownerID string) bool

	ShowNearBy(location *models.Location)(string,error)
}
