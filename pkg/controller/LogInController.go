package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vds/RestaurantManagement/pkg/database"
	"github.com/vds/RestaurantManagement/pkg/encryption"
	"github.com/vds/RestaurantManagement/pkg/middleware"
	"github.com/vds/RestaurantManagement/pkg/models"
	"log"
	"net/http"
)


type LogInController struct{
	database.Database
}

func NewLogInController(db database.Database)*LogInController{
	lc:=new(LogInController)
	lc.Database=db
	return lc
}
func(l *LogInController)LogIn(c *gin.Context){
	userType:=c.Param("userType")
	isValid:=middleware.IsValidUserType(userType)
	if !isValid{
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	key:=""
	switch userType{
	case middleware.Admin:
		key=middleware.AdminKey
	case middleware.SuperAdmin:
		key=middleware.SuperAdminKey
	default:
		key=middleware.OwnerKey
	}
	var cred models.Credentials
	err:=c.ShouldBindJSON(&cred)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	email,err:=l.LogInUser(userType,&cred)
	if err!=nil{
		fmt.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	token,err:=encryption.CreateToken(email,key)
	if err!=nil{
		log.Printf("%v",err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Header().Set("token", token)
	c.JSON(http.StatusOK,gin.H{
		"msg":"Login Successful",
	})
}