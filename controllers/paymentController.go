package controllers

import (
	"dodolan/models"
	"dodolan/utils/token"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type paymentDTO struct {
	Amount   float64
	OrderUid string `json:"uid_order"`
}

// GetPayments godoc
// @Summary Get all payments.
// @Description Get a list of payments.
// @Tags Payments
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /payments [get]
func GetPayments(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var payments []models.Payment
	db.Find(&payments)

	userToken, _ := token.ExtractTokenRole(c)
	//harus role cusstomer
	if userToken != "CUSTOMER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are unauthorized to access this resource, this resource for customer user",
		})
		return
	}

	if len(payments) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data Kosong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": payments})
}

// GetPaymentsById godoc
// @Summary Get payment by id
// @Description Get a payments by id.
// @Tags Payments
// @Produce json
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param id path string true "payments id"
// @Success 200 {object} map[string]interface{}
// @Router /payments/{id} [get]
func GetPaymentById(c *gin.Context) {
	var payment models.Payment

	db := c.MustGet("db").(*gorm.DB)

	userToken, _ := token.ExtractTokenRole(c)
	//harus role cusstomer
	if userToken != "CUSTOMER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are unauthorized to access this resource, this resource for customer user",
		})
		return
	}

	if err := db.Where("payment_id = ?", c.Param("id")).First(&payment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "payment not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": payment})
}

// GetOrdersByPaymentId godoc
// @Summary Get Orders By Payment ID.
// @Description Get all Orders by payments Id.
// @Tags Payments
// @Produce json
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param id path string true "payment id"
// @Success 200 {object} map[string]interface{}
// @Router /payments/{id}/orders [get]
func GetOrdersByPaymentId(c *gin.Context) {
	var orders []models.Order

	db := c.MustGet("db").(*gorm.DB)

	userToken, _ := token.ExtractTokenRole(c)
	//harus role cusstomer
	if userToken != "CUSTOMER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are unauthorized to access this resource, this resource for customer user",
		})
		return
	}

	if err := db.Where("payment_id = ?", c.Param("id")).First(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": orders})
}

// CreatePayment godoc
// @Summary Create New payment.
// @Description Creating a new payment.
// @Tags Payments
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body paymentDTO true "the body to create a new payments"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /payments [post]
func CreatePayment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input paymentDTO
	// var payments models.payments
	var order models.Order

	//harus bentuk json
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

	fmt.Println("cek nilai input order uid ", input.OrderUid)
	if err := db.Where("uid_order = ?", input.OrderUid).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record not found!"})
		return
	}
	fmt.Println("cek nilai order id ", order.OrderId)
	newPaymentUid := generateRandomID()
	//create Createpayments
	payments := models.Payment{
		UidPayment:  newPaymentUid,
		PaymentDate: time.Now(),
		Amount:      input.Amount,
		OrderId:     order.OrderId,
	}
	db.Create(&payments)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditambah", "data": payments})
}

// UpdatePayment godoc
// @Summary Update payments.
// @Description Update payments by id.
// @Tags Payments
// @Produce json
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param id path string true "payments id"
// @Param Body body paymentDTO true "the body to update an payments"
// @Success 200 {object} map[string]interface{}
// @Router /payments/{id} [put]
func UpdatePayment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var payments models.Payment
	//cek dulu uid payments yg ingin diupdate
	if err := db.Where("uid_payment = ?", c.Param("id")).First(&payments).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record not found!"})
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

	//harus bentuk json
	var input paymentDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	var updatedInput models.Payment
	updatedInput.Amount = input.Amount
	updatedInput.UpdatedAt = time.Now()

	db.Model(&payments).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Berhasil update payments", "data": payments})
}

// DeletePayment godoc
// @Summary Delete one payment.
// @Description Delete a payment by id.
// @Tags Payments
// @Produce json
// @Param id path string true "payment id"
// @Success 200 {object} map[string]interface{}
// @Router /payments/{id} [delete]
func DeletePayment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var payments models.Payment

	userToken, _ := token.ExtractTokenRole(c)
	//harus role cusstomer
	if userToken != "CUSTOMER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are unauthorized to access this resource, this resource for customer user",
		})
		return
	}

	//cek payments ada atau tidak
	if err := db.Where("uid_payment = ?", c.Param("id")).First(&payments).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "payments not found!"})
		return
	}
	db.Delete(&payments)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil dihapus", "data": payments})
}
