package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vds/RestaurantManagement/pkg/database"
	"github.com/vds/RestaurantManagement/pkg/middleware"
	"github.com/vds/RestaurantManagement/pkg/models"
	"net/http"
)


type RegisterController struct{
	database.Database
}

func NewRegisterController(db database.Database) *RegisterController{
	regController:=new(RegisterController)
	regController.Database=db
	return regController
}

func (r RegisterController)Register(c *gin.Context){
	userType:=c.Param("userType")
	if userType!=middleware.Admin && userType!=middleware.SuperAdmin{
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	var user models.User
	err:=c.ShouldBindJSON(&user)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err=r.CreateUser(userType,&user)
	if err!=nil{
		if err==database.ErrDupEmail{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"Registration Successful",
	})
}



