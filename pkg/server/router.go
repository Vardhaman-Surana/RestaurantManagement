package server

import (
	"github.com/gin-gonic/gin"
	"github.com/vds/RestaurantManagement/pkg/controller"
	"github.com/vds/RestaurantManagement/pkg/database"
	"github.com/vds/RestaurantManagement/pkg/middleware"
)
type Router struct{
	db database.Database
}

func NewRouter(db database.Database)(*Router,error){
	router := new(Router)
	router.db = db
	return router,nil
}
func (r *Router)Create(port string)error{
	//Controllers
	regController:=controller.NewRegisterController(r.db)
	logInController:=controller.NewLogInController(r.db)
	ownerController:=controller.NewOwnerController(r.db)
	resController:=controller.NewRestaurantController(r.db)
	menuController:=controller.NewMenuController(r.db)
	//Router and routes
	ginRouter:=gin.Default()

	ginRouter.POST("/register/:userType",regController.Register)
	ginRouter.POST("/login/:userType",logInController.LogIn)

	manage:=ginRouter.Group("/manage/:userType")
	manage.Use(middleware.AuthMiddleware)
	{
		manage.POST("/owners/add",ownerController.RegisterOwners)
		manage.DELETE("/owners/remove",ownerController.DeleteOwners)
		manage.GET("/owners",ownerController.GetOwners)
		// for restaurants
		manage.POST("/restaurants/add",resController.AddRestaurant)
		manage.GET("/restaurants",resController.GetRestaurants)
		manage.DELETE("/restaurants/remove",resController.DeleteRestaurants)
		manage.PUT("/restaurants/edit",resController.EditRestaurant)
		// for menu
		manage.POST("restaurants/menu/:resID/add",menuController.AddDishes)
		manage.DELETE("restaurants/menu/:resID/remove",menuController.DeleteDishes)
		manage.GET("restaurants/menu/:resID",menuController.GetMenu)
		manage.PUT("/restaurants/menu/:resID/edit",menuController.EditDishes)
	}
	owners:=ginRouter.Group("/owners")
	owners.Use(middleware.AuthMiddleware)
	{
		owners.GET("/restaurants",ownerController.GetRestaurants)
		owners.GET("/restaurants/menu/:resID",ownerController.GetMenu)
		owners.POST("/restaurants/menu/:resID/add",ownerController.AddDishes)
		owners.DELETE("/restaurants/menu/:resID/remove",ownerController.DeleteDishes)
		owners.PUT("/restaurants/menu/:resID/edit",ownerController.EditDishes)
	}
	ginRouter.GET("/restaurantsNearBy",resController.GetNearBy)



	return ginRouter.Run(port)
}