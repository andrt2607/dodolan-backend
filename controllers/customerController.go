package controllers

import (
	"dodolan/models"
	"dodolan/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type customerDTO struct {
	FirstName  string `json:"first_name" validate:"required, min=4,max=15"`
	LastName   string `json:"last_name" validate:"required, min=4,max=15"`
	Email      string `json:"email" validate:"required, email"`
	Username   string `json:"username" validate:"required, min=4,max=16"`
	Password   string `json:"password" validate:"required max=8"`
	Address    string `json:"address" validate:"required"`
	City       string `json:"city" validate:"required"`
	Country    string `json:"country" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required, len=5"`
	Phone      string `json:"phone" validate:"required, numeric, len=12"`
}

type updateCustomerDTO struct {
	FirstName  string `json:"first_name" validate:"required, min=4,max=15"`
	LastName   string `json:"last_name" validate:"required, min=4,max=15"`
	Email      string `json:"email" validate:"required, email"`
	Address    string `json:"address" validate:"required"`
	City       string `json:"city" validate:"required"`
	Country    string `json:"country" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required, len=5"`
	Phone      string `json:"phone" validate:"required, numeric, len=12"`
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
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data Customer berhasil ditemukan", "data": customers})
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
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Customer berhasil ditemukan", "data": Customer})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Products not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Products berhasil ditemukan", "data": products})
}

// CreateCustomer godoc
// @Summary Register New Customer.
// @Description Creating a new Customer.
// @Tags Customer
// @Param Body body customerDTO true "the body to create a new Customer"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /customers/register [post]
func CreateCustomer(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input customerDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	Customer := models.Customer{
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Email:      input.Email,
		Username:   input.Username,
		Password:   hashedPassword,
		Address:    input.Address,
		City:       input.City,
		Country:    input.Country,
		PostalCode: input.PostalCode,
		Phone:      input.Phone,
	}
	validate := validator.New()
	errValidate := validate.Struct(Customer)
	if errValidate != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": errValidate.Error(),
		})
		return
	}
	db.Create(&Customer)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Customer berhasil ditambah", "data": Customer})
}

// UpdateCustomer godoc
// @Summary Update Customer.
// @Description Update Customer by id.
// @Tags Customer
// @Produce json
// @Param id path string true "Customer id"
// @Param Body body updateCustomerDTO true "the body to update an Customer"
// @Success 200 {object} map[string]interface{}
// @Router /customers/{id} [put]
func UpdateCustomer(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var Customer models.Customer
	if err := db.Where("customer_id = ?", c.Param("id")).First(&Customer).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Customer not found!"})
		return
	}
	var input updateCustomerDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	var updatedInput models.Customer
	updatedInput.FirstName = input.FirstName
	updatedInput.LastName = input.LastName
	updatedInput.Address = input.Address
	updatedInput.City = input.City
	updatedInput.Country = input.Country
	updatedInput.PostalCode = input.PostalCode
	updatedInput.Phone = input.Phone
	updatedInput.UpdatedAt = time.Now()

	validate := validator.New()
	errValidate := validate.Struct(Customer)
	if errValidate != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": errValidate.Error(),
		})
		return
	}

	db.Model(&Customer).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Berhasil update Customer " + c.Param("id")})
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

	if err := db.Where("customer_id = ?", c.Param("id")).First(&Customer).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Customer not found!"})
		return
	}
	db.Delete(&Customer)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data customer id " + c.Param("id") + "berhasil dihapus"})
}
