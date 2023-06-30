package controllers

import (
	"dodolan/models"
	"dodolan/utils/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type productDTO struct {
	Name        string
	Description string
	Price       float64
	Stock       int64
	SellerId    int `json:"seller_id"`
}

// GetProducts godoc
// @Summary Get all Product.
// @Description Get a list of Product.
// @Tags Product
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /products [get]
func GetProducts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var products []models.Product
	db.Find(&products)

	//jika data kosong
	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data Kosong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": products})
}

// GetProductById godoc
// @Summary Get Product.
// @Description Get a Product by id.
// @Tags Product
// @Produce json
// @Param id path string true "product id"
// @Success 200 {object} models.Product
// @Router /products/{id} [get]
func GetProductById(c *gin.Context) {
	var product models.Product

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("product_id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Product not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": product})
}

// CreateProduct godoc
// @Summary Create New Product.
// @Description Creating a new Product.
// @Tags Product
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body productDTO true "the body to create a new product"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input productDTO
	var seller models.Seller

	//harus bentuk json
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	userToken, _ := token.ExtractTokenRole(c)
	if userToken != "SELLER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are unauthorized to access this resource, this resource for SELLER user",
		})
		return
	}

	//cari id seller dulu
	if err := db.Where("seller_id = ?", input.SellerId).First(&seller).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "seller not found"})
		return
	}

	//create CreateProduct
	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		SellerId:    input.SellerId,
	}
	db.Create(&product)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditambah", "data": product})
}

// GetSellerByProduct godoc
// @Summary Get All Seller by product id.
// @Description Get all seller by product id .
// @Tags Product
// @Produce json
// @Param id path string true "Product Id"
// @Success 200 {object} []models.Seller
// @Router /products/{id}/sellers [get]
func GetSellerByProduct(c *gin.Context) { // Get model if exist
	var sellers []models.Seller

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("seller_id = ?", c.Param("id")).Find(&sellers).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sellers})
}

// UpdateProduct godoc
// @Summary Update Product.
// @Description Update Product by id.
// @Tags Product
// @Produce json
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param id path string true "seller id"
// @Param Body body productDTO true "the body to update an product"
// @Success 200 {object} map[string]interface{}
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var product models.Product
	var seller models.Seller
	//cek dulu id product yg ingin diupdate
	if err := db.Where("product_id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record not found!"})
		return
	}

	userToken, _ := token.ExtractTokenRole(c)
	if userToken != "SELLER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are unauthorized to access this resource, this resource for SELLER user",
		})
		return
	}

	//harus bentuk json
	var input productDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	//cek seller ada atau tidak
	if err := db.Where("seller_id = ?", input.SellerId).First(&seller).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Seller not found!"})
		return
	}

	var updatedInput models.Product
	updatedInput.Name = input.Name
	updatedInput.Description = input.Description
	updatedInput.SellerId = input.SellerId
	updatedInput.UpdatedAt = time.Now()

	db.Model(&product).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Berhasil update product", "data": product})
}

// DeleteProduct godoc
// @Summary Delete one product.
// @Description Delete a product by id.
// @Tags Product
// @Produce json
// @Security BearerToken
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param id path string true "product id"
// @Success 200 {object} map[string]interface{}
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var product models.Product

	userToken, _ := token.ExtractTokenRole(c)
	if userToken != "SELLER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are unauthorized to access this resource, this resource for SELLER user",
		})
		return
	}

	//cek seller ada atau tidak
	if err := db.Where("product_id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Seller not found!"})
		return
	}
	db.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil dihapus", "data": product})
}
