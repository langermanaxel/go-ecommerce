package routes

import (
	"go-ecommerce/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incoming_routes *gin.Engine) {
	incoming_routes.POST("/users/signup", controllers.Signup())
	incoming_routes.POST("/users/login", controllers.Login())
	incoming_routes.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	incoming_routes.GET("/users/productview", controllers.SearchProduct())
	incoming_routes.GET("/users/search", controllers.SearchProductByQuery())
}
