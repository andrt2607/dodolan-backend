package routes

import (
	"dodolan/controllers"
	"dodolan/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	r.POST("/customer/login", controllers.LoginCustomer)
	r.POST("/seller/login", controllers.LoginSeller)

	r.GET("/products", controllers.GetProducts)
	r.GET("/products/:id", controllers.GetProductById)
	r.GET("/products/:id/sellers", controllers.GetProductsBySellerId)
	productsMiddlewareRoute := r.Group("/products")
	productsMiddlewareRoute.Use(middleware.JwtAuthMiddleware())
	productsMiddlewareRoute.POST("/", controllers.CreateProduct)
	productsMiddlewareRoute.PUT("/:id", controllers.UpdateProduct)
	productsMiddlewareRoute.DELETE("/:id", controllers.DeleteProduct)

	r.GET("/sellers", controllers.GetSellers)
	r.POST("/sellers/register", controllers.CreateSeller)
	r.GET("/sellers/:id", controllers.GetSellerById)
	r.PUT("/sellers/:id", controllers.UpdateSeller)
	r.DELETE("sellers/:id", controllers.DeleteSeller)

	r.GET("/customers", controllers.GetCustomers)
	r.POST("/customers/register", controllers.CreateCustomer)
	r.GET("/customers/:id", controllers.GetCustomerById)
	r.PUT("/customers/:id", controllers.UpdateCustomer)
	r.DELETE("customers/:id", controllers.DeleteCustomer)

	ordersMiddlewareRoute := r.Group("/orders")
	ordersMiddlewareRoute.Use(middleware.JwtAuthMiddleware())
	ordersMiddlewareRoute.GET("/items", controllers.GetOrderItems)
	ordersMiddlewareRoute.GET("/", controllers.GetOrders)
	ordersMiddlewareRoute.POST("/", controllers.CreateOrder)
	ordersMiddlewareRoute.PUT("/:id", controllers.UpdateOrderItem)
	ordersMiddlewareRoute.DELETE("/:id", controllers.DeleteOrder)

	paymentsMiddlewareRoute := r.Group("/payments")
	paymentsMiddlewareRoute.Use(middleware.JwtAuthMiddleware())
	paymentsMiddlewareRoute.GET("/", controllers.GetPayments)
	paymentsMiddlewareRoute.GET("/:id", controllers.GetPaymentById)
	paymentsMiddlewareRoute.POST("/", controllers.CreatePayment)
	paymentsMiddlewareRoute.PUT("/:id", controllers.UpdatePayment)
	paymentsMiddlewareRoute.DELETE("/:id", controllers.DeletePayment)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
