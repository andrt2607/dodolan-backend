package controllers

import (
	"dodolan/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type customerDTO struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
}

// GetCustomers godoc
// @Summary Get all Customer.
// @Description Get a list of Customer.
// @Tags Customer
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /customers [get]
func GetCustomers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var customers []models.Customer
	db.Find(&customers)

	if len(customers) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data Kosong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": customers})
}

// GetCustomerById godoc
// @Summary Get Customer.
// @Description Get a Customer by id.
// @Tags Customer
// @Produce json
// @Param id path string true "Customer id"
// @Success 200 {object} map[string]interface{}
// @Router /customers/{id} [get]
func GetCustomerById(c *gin.Context) {
	var Customer models.Customer

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("customer_id = ?", c.Param("id")).First(&Customer).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Customer not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": Customer})
}

// GetProductsByCustomerId godoc
// @Summary Get Products.
// @Description Get all Products by Customer Id.
// @Tags Customer
// @Produce json
// @Param id path string true "Customer id"
// @Success 200 {object} map[string]interface{}
// @Router /customers/{id}/products [get]
func GetProductsByCustomerId(c *gin.Context) {
	var products []models.Product

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("customer_id = ?", c.Param("id")).First(&products).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": products})
}

// CreateCustomer godoc
// @Summary Create New Customer.
// @Description Creating a new Customer.
// @Tags Customer
// @Param Body body input true "the body to create a new Customer"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /customers [post]
func CreateCustomer(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input customerDTO
	// var Customer models.Customer

	//harus bentuk json
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	//create CreateCustomer
	Customer := models.Customer{
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Email:      input.Email,
		Username:   input.Username,
		Password:   input.Password,
		Address:    input.Address,
		City:       input.City,
		Country:    input.Country,
		PostalCode: input.PostalCode,
		Phone:      input.Phone,
		// 	FirstName  string
		// LastName   string
		// Email      string
		// Username   string
		// Password   string
		// Address    string
		// City       string
		// Country    string
		// PostalCode string
		// Phone      string
	}
	db.Create(&Customer)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditambah", "data": Customer})
}

// UpdateCustomer godoc
// @Summary Update Customer.
// @Description Update Customer by id.
// @Tags Customer
// @Produce json
// @Param id path string true "Customer id"
// @Param Body body input true "the body to update an Customer"
// @Success 200 {object} map[string]interface{}
// @Router /customers/{id} [put]
func UpdateCustomer(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var Customer models.Customer
	// var Customer models.Customer
	//cek dulu id Customer yg ingin diupdate
	if err := db.Where("customer_id = ?", c.Param("id")).First(&Customer).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record not found!"})
		return
	}
	//harus bentuk json
	var input customerDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	//cek Customer ada atau tidak
	// if err := db.Where("id = ?", input.CustomerId).First(&Customer).Error; err != nil {
	//       c.JSON(http.StatusBadRequest, gin.H{"error":true, "message": "Customer not found!"})
	//       return
	//   }

	var updatedInput models.Customer
	updatedInput.FirstName = input.FirstName
	updatedInput.LastName = input.LastName
	updatedInput.Username = input.Username
	updatedInput.Password = input.Password
	updatedInput.Address = input.Address
	updatedInput.City = input.City
	updatedInput.Country = input.Country
	updatedInput.PostalCode = input.PostalCode
	updatedInput.Phone = input.Phone
	// FirstName  string
	// LastName   string
	// Email      string
	// Username   string
	// Password   string
	// Address    string
	// City       string
	// Country    string
	// PostalCode string
	// Phone      string
	// updatedInput.CustomerId = input.CustomerId
	updatedInput.UpdatedAt = time.Now()

	db.Model(&Customer).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Berhasil update Customer", "data": Customer})
}

// DeleteCustomer godoc
// @Summary Delete one Customer.
// @Description Delete a Customer by id.
// @Tags Customer
// @Produce json
// @Param id path string true "Customer id"
// @Success 200 {object} map[string]interface{}
// @Router /customers/{id} [delete]
func DeleteCustomer(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var Customer models.Customer

	//cek Customer ada atau tidak
	if err := db.Where("customer_id = ?", c.Param("id")).First(&Customer).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Customer not found!"})
		return
	}
	db.Delete(&Customer)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil dihapus", "data": Customer})
}
