package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/vds/RestaurantManagement/pkg/database"
	"github.com/vds/RestaurantManagement/pkg/middleware"
	"github.com/vds/RestaurantManagement/pkg/models"
	"net/http"
)

type RestaurantController struct{
	database.Database
}

func NewRestaurantController(db database.Database) *RestaurantController{
	resController:=new(RestaurantController)
	resController.Database=db
	return resController
}

func(r *RestaurantController)GetNearBy(c *gin.Context){
	var location models.Location
	err:=c.ShouldBindJSON(&location)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	jsonData:=[]struct{
		Name string `json"name" binding"required"`
	}{}
	stringData,err:=r.ShowNearBy(&location)
	if stringData==""{
		c.JSON(http.StatusOK, gin.H{
			"msg":"No restaurants to show",
		})
		return
	}
	err=json.Unmarshal([]byte(stringData),&jsonData)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,jsonData)

}

func(r *RestaurantController)GetRestaurants(c *gin.Context){
	creatorID,_:=c.Get("email")
	jsonData:=&[]models.RestaurantOutput{}
	userType:=c.Param("userType")
	if userType!=middleware.Admin && userType!=middleware.SuperAdmin{
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	var stringData string
	var err error
	if userType==middleware.SuperAdmin{
		stringData,err=r.ShowRestaurantsForSuper()
		if err!=nil{
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}else {
		stringData, err = r.ShowRestaurants(creatorID.(string))
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
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
func(r *RestaurantController)AddRestaurant(c *gin.Context){
	creatorID,_:=c.Get("email")
	var restaurant models.Restaurant
	err:=c.ShouldBindJSON(&restaurant)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	restaurant.CreatorID=creatorID.(string)
	if restaurant.OwnerID!=""{
		isExisting:=r.IsExistingOwner(restaurant.OwnerID)
		if !isExisting{
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Owner does not exist please create the owner or try creating without owner",
			})
			return
		}
	}
	err=r.InsertRestaurant(&restaurant)
	if err!=nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"Restaurant added",
	})
}
func(r *RestaurantController)DeleteRestaurants(c *gin.Context){
	userType:=c.Param("userType")
	creatorID,_:=c.Get("email")
	var ids []models.ResID
	err:=c.ShouldBindJSON(&ids)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if userType==middleware.SuperAdmin{
		err=r.RemoveRestaurantsBySuper(ids)
		if err!=nil{
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "Restaurants Removed Successfully",
		})
		return
	}else if userType==middleware.Admin{
		err=r.CheckRestaurantCreator(creatorID.(string),ids)
		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{
				"msg":"Can't Delete restaurant created by others",
			})
			return
		}

		err=r.RemoveRestaurants(creatorID.(string),ids)
		if err!=nil{
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "Restaurants Removed Successfully",
		})
		return
	}else{
		c.Writer.WriteHeader(http.StatusNotFound)
	}
}
func(r *RestaurantController)EditRestaurant(c *gin.Context){
	userType:=c.Param("userType")
	creatorID,_:=c.Get("email")
	var restaurant models.RestaurantOutput
	err:=c.ShouldBindJSON(&restaurant)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	id:=[]models.ResID{
		{ID:restaurant.ID},
	}
	if userType!=middleware.SuperAdmin {
		err := r.CheckRestaurantCreator(creatorID.(string), id)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "Can't update restaurant created by others",
			})
			return
		}
	}

	err=r.UpdateRestaurant(&restaurant)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "Restaurant Updated Successfully",
	})
}