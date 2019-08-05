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

type MenuController struct{
	database.Database
}

func NewMenuController(db database.Database) *MenuController{
	menuController:=new(MenuController)
	menuController.Database=db
	return menuController
}
func (m *MenuController)GetMenu(c *gin.Context){
	userType:=c.Param("userType")
	creatorID,_:=c.Get("email")
	value:=c.Param("resID")
	resID,_:=strconv.Atoi(value)
	jsonData:=&[]models.DishOutput{}
	id:=[]models.ResID{
		{ID:resID},
	}
	if userType!=middleware.SuperAdmin {
		err := m.CheckRestaurantCreator(creatorID.(string), id)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "Can't Show the menu restaurant created by others",
			})
			return
		}
	}
	var stringData string
	stringData,err:=m.ShowMenu(resID)
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

func(m *MenuController)AddDishes(c *gin.Context){
	userType:=c.Param("userType")
	value:=c.Param("resID")
	creatorID,_:=c.Get("email")
	resID,_:=strconv.Atoi(value)
	id:=[]models.ResID{
		{ID:resID},
	}
	if userType!=middleware.SuperAdmin {
		err := m.CheckRestaurantCreator(creatorID.(string), id)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "Can't Add dishes restaurant created by others",
			})
			return
		}
	}
	var dishes []models.Dish
	err:=c.ShouldBindJSON(&dishes)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err=m.InsertDishes(dishes,resID)
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

func (m *MenuController)DeleteDishes(c *gin.Context){
	userType:=c.Param("userType")
	value:=c.Param("resID")
	creatorID,_:=c.Get("email")
	resID,_:=strconv.Atoi(value)
	id:=[]models.ResID{
		{ID:resID},
	}
	if userType==middleware.Admin{
		err:=m.CheckRestaurantCreator(creatorID.(string),id)
		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{
				"msg":"Can't Remove dishes of restaurant created by others",
			})
			return
		}
	}
	var dishIDs []models.DishID
	err:=c.ShouldBindJSON(&dishIDs)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err=m.RemoveDishes(dishIDs)
	if err!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"Dishes Removed",
	})
}

func (m *MenuController)EditDishes(c *gin.Context){
	userType:=c.Param("userType")
	value:=c.Param("resID")
	creatorID,_:=c.Get("email")
	resID,_:=strconv.Atoi(value)
	id:=[]models.ResID{
		{ID:resID},
	}
	if userType==middleware.Admin{
		err:=m.CheckRestaurantCreator(creatorID.(string),id)
		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{
				"msg":"Can't Remove dishes of restaurant created by others",
			})
			return
		}
	}
	var dishes []models.DishOutput
	err:=c.ShouldBindJSON(&dishes)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err=m.UpdateDishes(dishes)
	if err!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":"Dishes Updated successfully",
	})
}

