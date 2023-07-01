package controllers

import (
	"crypto/rand"
	"dodolan/models"
	"dodolan/utils/token"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func generateRandomID() string {
	length := 32
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		panic(err)
	}
	randomID := base64.URLEncoding.EncodeToString(buffer)
	return randomID
}

type orderDTO struct {
	ProductId   int `json:"product_id"`
	Quantity    int64
	Price       float64
	TotalAmount float64 `json:"total_amount"`
	CustomerId  int     `json:"customer_id"`
}

// GetOrder godoc
// @Summary Get all Order.
// @Description Get a list of Order.
// @Tags Order
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /order [get]
func GetOrders(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var orders []models.Order
	db.Find(&orders)

	userToken, _ := token.ExtractTokenRole(c)
	//harus role cusstomer
	if userToken != "CUSTOMER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are unauthorized to access this resource, this resource for customer user",
		})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data Kosong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data order id " + c.Param("id") + " berhasil ditemukan", "data": orders})
}

// GetOrderItem godoc
// @Summary Get all OrderItems.
// @Description Get a list of OrderItem.
// @Tags Order
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /order_item [get]
func GetOrderItems(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var orderItems []models.OrderItem
	db.Find(&orderItems)

	userToken, _ := token.ExtractTokenRole(c)
	//harus role cusstomer
	if userToken != "CUSTOMER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are unauthorized to access this resource, this resource for customer user",
		})
		return
	}

	if len(orderItems) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data Kosong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data orderitem id " + c.Param("id") + "berhasil ditemukan", "data": orderItems})
}

// CreateOrder godoc
// @Summary Create Order.
// @Description Creating a new Order.
// @Tags Order
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body orderDTO true "the body to create a new Order"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input orderDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	userToken, _ := token.ExtractTokenRole(c)
	//harus role cusstomer
	if userToken != "CUSTOMER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are unauthorized to access this resource, this resource for customer user",
		})
		return
	}

	newOrderId := generateRandomID()
	newOrderItemId := generateRandomID()
	//create CreateSeller
	order := models.Order{
		UidOrder:    newOrderId,
		OrderDate:   time.Now(),
		CustomerId:  input.CustomerId,
		TotalAmount: input.TotalAmount,
	}
	validate := validator.New()
	errValidate := validate.Struct(order)
	if errValidate != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": errValidate.Error(),
		})
		return
	}
	db.Create(&order)
	var targetOrder models.Order

	// db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("uid_order = ?", newOrderId).First(&targetOrder).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Order not found!"})
		return
	}
	orderItem := models.OrderItem{
		UidOrderItem: newOrderItemId,
		OrderId:      targetOrder.OrderId,
		Quantity:     input.Quantity,
		ProductId:    input.ProductId,
	}

	db.Create(&orderItem)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditambah"})
}

// UpdateOrder godoc
// @Summary Update Order.
// @Description Update Order by id.
// @Tags Order
// @Produce json
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param id path string true "order id"
// @Param Body body orderDTO true "the body to update an Order"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id} [put]
func UpdateOrderItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var orderItem models.OrderItem
	var order models.Order
	userToken, _ := token.ExtractTokenRole(c)
	//harus role cusstomer
	if userToken != "CUSTOMER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are unauthorized to access this resource, this resource for customer user",
		})
		return
	}
	//where ini untuk concat query update order dibawah
	if err := db.Where("uid_order_item = ?", c.Param("id")).First(&orderItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record order item not found!"})
		return
	}
	fmt.Println("value 1 : ", orderItem.UidOrderItem)
	//where ini untuk concat query update order dibawah
	if err := db.Where("order_id = ?", orderItem.OrderId).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record order not found!"})
		return
	}
	var input orderDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	var targetOrder models.Order
	if err := db.Where("order_id = ?", orderItem.OrderId).First(&targetOrder).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Order not found!"})
		return
	}
	var updatedInput models.OrderItem
	updatedInput.ProductId = input.ProductId
	updatedInput.Quantity = input.Quantity
	updatedInput.Order.TotalAmount = input.TotalAmount
	updatedInput.Order.UpdatedAt = time.Now()
	var updatedOrderInput models.Order
	updatedOrderInput.TotalAmount = input.TotalAmount
	updatedOrderInput.UpdatedAt = time.Now()
	updatedInput.UpdatedAt = time.Now()

	validate := validator.New()
	errValidate := validate.Var(input.TotalAmount, "required")
	if errValidate != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": errValidate.Error(),
		})
		return
	}
	db.Model(&orderItem).Updates(updatedInput)

	db.Model(&order).Updates(updatedOrderInput)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Berhasil update order id " + c.Param("id") + "item", "data": updatedInput})
}

// DeleteOrder godoc
// @Summary Delete one Order.
// @Description Delete a Order by id.
// @Tags Order
// @Produce json
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param id path string true "Order id"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id} [delete]
func DeleteOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var orderItem models.OrderItem
	var order models.Order
	//
	userToken, _ := token.ExtractTokenRole(c)
	//harus role cusstomer
	if userToken != "CUSTOMER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are unauthorized to access this resource, this resource for customer user",
		})
		return
	}
	//where ini untuk concat query update order dibawah
	if err := db.Where("uid_order_item = ?", c.Param("id")).First(&orderItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record order item not found!"})
		return
	}
	//where ini untuk concat query update order dibawah
	if err := db.Where("order_id = ?", orderItem.OrderId).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record order not found!"})
		return
	}
	db.Delete(&orderItem)
	db.Delete(&order)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data order id " + c.Param("id") + "berhasil dihapus", "data": orderItem})
}
