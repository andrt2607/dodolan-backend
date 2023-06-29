package routes

import (
	"dodolan/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// set db to gin context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})
	r.GET("/products", controllers.GetProducts)
	r.POST("/products", controllers.CreateProduct)
	r.GET("/products/:id", controllers.GetProductById)
	r.GET("/products/:id/sellers", controllers.GetProductsBySellerId)
	r.PUT("/products/:id", controllers.UpdateProduct)
	r.DELETE("products/:id", controllers.DeleteProduct)

	r.GET("/sellers", controllers.GetSellers)
	r.POST("/sellers", controllers.CreateSeller)
	r.GET("/sellers/:id", controllers.GetSellerById)
	r.PUT("/sellers/:id", controllers.UpdateSeller)
	r.DELETE("sellers/:id", controllers.DeleteSeller)

	r.GET("/customers", controllers.GetCustomers)
	r.POST("/customers", controllers.CreateCustomer)
	r.GET("/customers/:id", controllers.GetCustomerById)
	r.PUT("/customers/:id", controllers.UpdateCustomer)
	r.DELETE("customers/:id", controllers.DeleteCustomer)

	// r.GET("/customers", controllers.GetCustomers)
	r.POST("/orders", controllers.CreateOrder)
	// r.GET("/customers/:id", controllers.GetCustomerById)
	r.PUT("/orders/:id", controllers.UpdateOrderItem)
	// r.DELETE("customers/:id", controllers.DeleteCustomer)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
