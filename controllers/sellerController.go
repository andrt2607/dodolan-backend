package controllers

import (
	"dodolan/models"
	"dodolan/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SellerDTO struct {
	Name     string
	Username string
	Password string
	Address  string
	Phone    string
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
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": sellers})
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
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": Seller})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditemukan", "data": products})
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
	// var seller models.Seller

	//harus bentuk json
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

	//create CreateSeller
	seller := models.Seller{
		Name:     input.Name,
		Username: input.Username,
		Password: hashedPassword,
		Address:  input.Address,
		Phone:    input.Phone,
	}
	db.Create(&seller)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil ditambah", "data": seller})
}

// UpdateSeller godoc
// @Summary Update Seller.
// @Description Update Seller by id.
// @Tags Seller
// @Produce json
// @Param id path string true "seller id"
// @Param Body body SellerDTO true "the body to update an Seller"
// @Success 200 {object} map[string]interface{}
// @Router /Sellers/{id} [put]
func UpdateSeller(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var Seller models.Seller
	// var seller models.Seller
	//cek dulu id Seller yg ingin diupdate
	if err := db.Where("seller_id = ?", c.Param("id")).First(&Seller).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Record not found!"})
		return
	}
	//harus bentuk json
	var input SellerDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	//cek seller ada atau tidak
	// if err := db.Where("id = ?", input.SellerId).First(&seller).Error; err != nil {
	//       c.JSON(http.StatusBadRequest, gin.H{"error":true, "message": "Seller not found!"})
	//       return
	//   }

	var updatedInput models.Seller
	updatedInput.Name = input.Name
	updatedInput.Address = input.Address
	// updatedInput.SellerId = input.SellerId
	updatedInput.UpdatedAt = time.Now()

	db.Model(&Seller).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Berhasil update Seller", "data": Seller})
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

	//cek seller ada atau tidak
	if err := db.Where("seller_id = ?", c.Param("id")).First(&Seller).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Seller not found!"})
		return
	}
	db.Delete(&Seller)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Data berhasil dihapus", "data": Seller})
}
