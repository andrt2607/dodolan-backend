package controllers

import (
	"crypto/rand"
	"dodolan/models"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func generateRandomID() string {
	// Menentukan panjang byte untuk ID acak
	// Disarankan untuk menggunakan setidaknya 16 byte (128 bit) untuk kekuatan keamanan yang baik
	// Di sini, kita menggunakan 32 byte (256 bit)
	length := 32

	// Membuat buffer dengan panjang yang ditentukan
	buffer := make([]byte, length)

	// Membaca byte acak ke dalam buffer
	_, err := rand.Read(buffer)
	if err != nil {
		panic(err)
	}

	// Mengonversi byte menjadi string dengan encoding base64
	// Ini menghasilkan string yang lebih panjang, tetapi memiliki karakter yang lebih aman untuk penggunaan di URL
	randomID := base64.URLEncoding.EncodeToString(buffer)

	// Mengembalikan ID acak
	return randomID
}

// import (
// 	"dodolan/models"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

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
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /order [get]
func GetOrders(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var orders []models.Order
	db.Find(&orders)

	if len(orders) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data Kosong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": orders})
}

// GetOrderItem godoc
// @Summary Get all OrderItems.
// @Description Get a list of OrderItem.
// @Tags Order
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /order_item [get]
func GetOrderItems(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var orderItems []models.OrderItem
	db.Find(&orderItems)

	if len(orderItems) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data Kosong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": orderItems})
}

// CreateOrder godoc
// @Summary Create Order.
// @Description Creating a new Order.
// @Tags Order
// @Param Body body input true "the body to create a new Order"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input orderDTO
	// var order models.Order
	// var orderItem models.OrderItem

	//harus bentuk json
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
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
// @Param id path string true "order id"
// @Param Body body input true "the body to update an Order"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id} [put]
func UpdateOrderItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var orderItem models.OrderItem
	var order models.Order
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
	// fmt.Println("ini value totalamount : ", input.TotalAmount)
	// fmt.Println("ini value totalamount dari order: ", targetOrder.TotalAmount)
	var updatedInput models.OrderItem
	updatedInput.ProductId = input.ProductId
	updatedInput.Quantity = input.Quantity
	updatedInput.Order.TotalAmount = input.TotalAmount
	updatedInput.Order.UpdatedAt = time.Now()
	var updatedOrderInput models.Order
	updatedOrderInput.TotalAmount = input.TotalAmount
	updatedOrderInput.UpdatedAt = time.Now()
	updatedInput.UpdatedAt = time.Now()

	db.Model(&orderItem).Updates(updatedInput)

	db.Model(&order).Updates(updatedOrderInput)
	fmt.Println("ini value id order item : ", orderItem.OrderId)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Berhasil update order item", "data": updatedInput})
}

// // DeleteSeller godoc
// // @Summary Delete one Seller.
// // @Description Delete a Seller by id.
// // @Tags Seller
// // @Produce json
// // @Param id path string true "Seller id"
// // @Success 200 {object} map[string]interface{}
// // @Router /Sellers/{id} [delete]
// func DeleteSeller(c *gin.Context) {
// 	db := c.MustGet("db").(*gorm.DB)
// 	var Seller models.Seller

// 	//cek seller ada atau tidak
// 	if err := db.Where("seller_id = ?", c.Param("id")).First(&Seller).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Seller not found!"})
// 		return
// 	}
// 	db.Delete(&Seller)
// 	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil dihapus", "data": Seller})
// }
