package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vds/RestaurantManagement/pkg/models"
	"net/http"
	"strings"
)
const Owner="owner"
const AdminKey="adminKey"
const SuperAdminKey="superAdminKey"
const OwnerKey="ownerKey"
const tokenExpireMessage="Token expired please login again"
const Admin="admin"
const SuperAdmin="superAdmin"

func AuthMiddleware(c *gin.Context){
	userType:=c.Param("userType")
	if strings.Contains(c.Request.URL.Path,"/manage") {
		isValid := IsValidUserType(userType)
		if !isValid {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
	key:=""
	switch userType{
	case Admin:
		key=AdminKey
	case SuperAdmin:
		key=SuperAdminKey
	default:
		key=OwnerKey
	}
	jwtKey:=[]byte(key)
	tokenStr:=c.Request.Header.Get("token")

	claims:=&models.Claims{}

	tkn,err:=jwt.ParseWithClaims(tokenStr,claims,func(token *jwt.Token)(interface{},error){
		return jwtKey,nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if strings.Contains(err.Error(), "expired") {
			fmt.Print(err)
			c.JSON(http.StatusUnauthorized,gin.H{
				"msg": tokenExpireMessage,
			})
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		fmt.Printf("%v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if !tkn.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("email",claims.Email)
	c.Next()
}

func IsValidUserType(userType string)bool{
	if userType!=Admin && userType!=SuperAdmin && userType!=Owner{
		return false
	}
	return true
}