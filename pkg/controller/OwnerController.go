package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/vds/RestaurantManagement/pkg/database"
	"github.com/vds/RestaurantManagement/pkg/middleware"
	"github.com/vds/RestaurantManagement/pkg/models"
	"net/http"
	"strconv"
)

type OwnerController struct{
	database.Database
}

func NewOwnerController(db database.Database)*OwnerController{
	ownerController:=new(OwnerController)
	ownerController.Database=db
	return ownerController
}

func(o *OwnerController)GetOwners(c *gin.Context){
	userType:=c.Param("userType")
	if userType!=middleware.Admin && userType!=middleware.SuperAdmin{
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	creatorID,_:=c.Get("email")
	jsonData:=&[]models.UserOutput{}
	var stringData string
	var err error
	if userType==middleware.SuperAdmin{
		stringData,err=o.ShowOwnersForSuper()
		if err!=nil{
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
	}else {
		stringData, err = o.ShowOwners(creatorID.(string))
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
	}
	err=json.Unmarshal([]byte(stringData),jsonData)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,jsonData)
}

func(o *OwnerController)RegisterOwners(c *gin.Context){
	creatorID,_:=c.Get("email")
	userType:=c.Param("userType")
	if userType!=middleware.Admin && userType!=middleware.SuperAdmin{
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	var user []models.User
	err:=c.ShouldBindJSON(&user)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err=o.CreateOwners(creatorID.(string),user)
	if err!=nil{
		if err!=database.ErrInternal{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"Owners created successfully",
	})
}

func(o *OwnerController)DeleteOwners(c *gin.Context){
	creatorID,_:=c.Get("email")
	userType:=c.Param("userType")
	if userType!=middleware.Admin && userType!=middleware.SuperAdmin{
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	var ownerIDs []models.UserID
	err:=c.ShouldBindJSON(&ownerIDs)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if userType==middleware.SuperAdmin{
		err= o.RemoveOwnersBySuper(ownerIDs)
		if err!=nil{
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"msg":"Owners removed successfully",
		})
		return
	}
	err=o.CheckOwnerCreator(creatorID.(string),ownerIDs)
	if err!=nil{
		c.JSON(http.StatusUnauthorized,gin.H{
			"msg":"Can't Delete owners created by others",
		})
		return
	}

	err=o.RemoveOwners(creatorID.(string),ownerIDs)
	if err!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"Owners removed successfully",
	})
}

func (o *OwnerController)GetRestaurants(c *gin.Context){
	ownerID,_:=c.Get("email")
	stringData,err:=o.GetOwnerRestaurants(ownerID.(string))
	if err!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData:=&[]models.RestaurantOutput{}
	if stringData==""{
		c.JSON(http.StatusOK, gin.H{
			"msg":"No restaurants to show",
		})
		return
	}
	err=json.Unmarshal([]byte(stringData),jsonData)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,jsonData)
}
func (o *OwnerController)GetMenu(c *gin.Context){
	ownerID,_:=c.Get("email")
	value:=c.Param("resID")
	resID,_:=strconv.Atoi(value)
	jsonData:=&[]models.DishOutput{}
	var stringData string
	err := o.CheckRestaurantOwner(ownerID.(string),resID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Can't Show the menu restaurant created by others",
		})
		return
	}
	stringData,err=o.ShowMenu(resID)
	if stringData==""{
		c.JSON(http.StatusOK, gin.H{
			"msg":"No Items to show",
		})
		return
	}
	err=json.Unmarshal([]byte(stringData),jsonData)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,jsonData)

}
func(o *OwnerController)AddDishes(c *gin.Context){
	ownerID,_:=c.Get("email")
	value:=c.Param("resID")
	resID,_:=strconv.Atoi(value)
	err := o.CheckRestaurantOwner(ownerID.(string),resID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Can't Add the dishes to restaurant created by others",
		})
		return
	}
	var dishes []models.Dish
	err=c.ShouldBindJSON(&dishes)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err=o.InsertDishes(dishes,resID)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":"Dishes Added to menu successfully",
	})
	return
}
func(o *OwnerController)DeleteDishes(c *gin.Context){
	ownerID,_:=c.Get("email")
	value:=c.Param("resID")
	resID,_:=strconv.Atoi(value)
	err := o.CheckRestaurantOwner(ownerID.(string),resID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Can't remove the dishes of restaurant created by others",
		})
		return
	}
	var dishIDs []models.DishID
	err=c.ShouldBindJSON(&dishIDs)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err=o.RemoveDishes(dishIDs)
	if err!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"Dishes Removed",
	})
}
func(o *OwnerController)EditDishes(c *gin.Context){
	ownerID,_:=c.Get("email")
	value:=c.Param("resID")
	resID,_:=strconv.Atoi(value)
	err := o.CheckRestaurantOwner(ownerID.(string),resID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Can't edit  dishes of restaurant created by others",
		})
		return
	}
	var dishes []models.DishOutput
	err=c.ShouldBindJSON(&dishes)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err=o.UpdateDishes(dishes)
	if err!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":"Dishes Updated successfully",
	})
}
