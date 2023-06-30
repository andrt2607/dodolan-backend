package routes

import (
	"dodolan/controllers"
	"dodolan/middleware"

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

	r.POST("/customer/login", controllers.LoginCustomer)
	r.POST("/seller/login", controllers.LoginSeller)

	//harus role seller
	r.GET("/products", controllers.GetProducts)
	r.GET("/products/:id", controllers.GetProductById)
	r.GET("/products/:id/sellers", controllers.GetProductsBySellerId)
	productsMiddlewareRoute := r.Group("/products")
	productsMiddlewareRoute.Use(middleware.JwtAuthMiddleware())
	// r.POST("/products", controllers.CreateProduct)
	// r.PUT("/products/:id", controllers.UpdateProduct)
	// r.DELETE("products/:id", controllers.DeleteProduct)
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

	// r.GET("/orders/items", controllers.GetOrderItems)
	// r.GET("/orders", controllers.GetOrders)
	ordersMiddlewareRoute := r.Group("/orders")
	ordersMiddlewareRoute.Use(middleware.JwtAuthMiddleware())
	ordersMiddlewareRoute.GET("/items", controllers.GetOrderItems)
	ordersMiddlewareRoute.GET("/", controllers.GetOrders)
	ordersMiddlewareRoute.POST("/", controllers.CreateOrder)
	ordersMiddlewareRoute.PUT("/:id", controllers.UpdateOrderItem)
	ordersMiddlewareRoute.DELETE("/:id", controllers.DeleteOrder)

	paymentsMiddlewareRoute := r.Group("/payments")
	paymentsMiddlewareRoute.Use(middleware.JwtAuthMiddleware())
	// paymentsMiddlewareRoute.GET("/items", controllers.GetOrderItems)
	paymentsMiddlewareRoute.GET("/", controllers.GetPayments)
	paymentsMiddlewareRoute.GET("/:id", controllers.GetPaymentById)
	paymentsMiddlewareRoute.POST("/", controllers.CreatePayment)
	paymentsMiddlewareRoute.PUT("/:id", controllers.UpdatePayment)
	paymentsMiddlewareRoute.DELETE("/:id", controllers.DeletePayment)
	// r.GET("/payments/", controllers.GetPayments)
	// r.GET("/payments/:id", controllers.GetPaymentById)
	// r.POST("/payments", controllers.CreatePayment)
	// r.PUT("/payments/:id", controllers.UpdatePayment)
	// r.DELETE("payments/:id", controllers.DeletePayment)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
