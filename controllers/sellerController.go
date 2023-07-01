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

type SellerDTO struct {
	Name     string `validate:"required, min=4, max=15"`
	Username string `validate:"required, min=4,max=16"`
	Password string `validate:"required max=8"`
	Address  string `validate:"required"`
	Phone    string `validate:"required, numeric, len=12"`
}

type SellerUpdateDTO struct {
	Name    string `validate:"required, min=4, max=15"`
	Address string `validate:"required"`
	Phone   string `validate:"required, numeric, len=12"`
}

// GetSeller godoc
// @Summary Get all Seller.
// @Description Get a list of Seller.
// @Tags Seller
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /sellers [get]
func GetSellers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var sellers []models.Seller
	db.Find(&sellers)

	if len(sellers) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data Kosong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data sellerberhasil ditemukan", "data": sellers})
}

// GetSellerById godoc
// @Summary Get Seller.
// @Description Get a Seller by id.
// @Tags Seller
// @Produce json
// @Param id path string true "Seller id"
// @Success 200 {object} map[string]interface{}
// @Router /sellers/{id} [get]
func GetSellerById(c *gin.Context) {
	var Seller models.Seller

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("seller_id = ?", c.Param("id")).First(&Seller).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Seller not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data seller id " + c.Param("id") + "berhasil ditemukan", "data": Seller})
}

// GetProductsBySellerId godoc
// @Summary Get Products by seller id.
// @Description Get all Products by Seller Id.
// @Tags Seller
// @Produce json
// @Param id path string true "Seller id"
// @Success 200 {object} map[string]interface{}
// @Router /sellers/{id}/products [get]
func GetProductsBySellerId(c *gin.Context) {
	var products []models.Product

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("seller_id = ?", c.Param("id")).First(&products).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Products not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data products berhasil ditemukan", "data": products})
}

// CreateSeller godoc
// @Summary Register New Seller.
// @Description Creating a new Seller.
// @Tags Seller
// @Param Body body SellerDTO true "the body to create a new Seller"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /sellers/register [post]
func CreateSeller(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input SellerDTO

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

	seller := models.Seller{
		Name:     input.Name,
		Username: input.Username,
		Password: hashedPassword,
		Address:  input.Address,
		Phone:    input.Phone,
	}
	validate := validator.New()
	errValidate := validate.Struct(seller)
	if errValidate != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": errValidate.Error(),
		})
		return
	}
	db.Create(&seller)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data seller berhasil ditambah", "data": seller})
}

// UpdateSeller godoc
// @Summary Update Seller.
// @Description Update Seller by id.
// @Tags Seller
// @Produce json
// @Param id path string true "seller id"
// @Param Body body SellerUpdateDTO true "the body to update an Seller"
// @Success 200 {object} map[string]interface{}
// @Router /Sellers/{id} [put]
func UpdateSeller(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var Seller models.Seller
	if err := db.Where("seller_id = ?", c.Param("id")).First(&Seller).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record not found!"})
		return
	}
	var input SellerUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	var updatedInput models.Seller
	updatedInput.Name = input.Name
	updatedInput.Address = input.Address
	updatedInput.UpdatedAt = time.Now()
	validate := validator.New()
	errValidate := validate.Struct(updatedInput)
	if errValidate != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": errValidate.Error(),
		})
		return
	}
	db.Model(&Seller).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Berhasil update Seller id " + c.Param("id"), "data": Seller})
}

// DeleteSeller godoc
// @Summary Delete one Seller.
// @Description Delete a Seller by id.
// @Tags Seller
// @Produce json
// @Param id path string true "Seller id"
// @Success 200 {object} map[string]interface{}
// @Router /Sellers/{id} [delete]
func DeleteSeller(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var Seller models.Seller

	if err := db.Where("seller_id = ?", c.Param("id")).First(&Seller).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Seller not found!"})
		return
	}
	db.Delete(&Seller)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data seller id " + c.Param("id") + "berhasil dihapus", "data": Seller})
}
