package encryption

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/vds/RestaurantManagement/pkg/models"
	"time"
)

func CreateToken(email string,key string) (string,error){
	jwtKey:=[]byte(key)
	expirationTime:=time.Now().Add(30*time.Minute)
	claims:=&models.Claims{
		Email:email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt:expirationTime.Unix(),
		},
	}
	//remember to change it later
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err!=nil{
		return "",err
	}
	return tokenString,nil
}
