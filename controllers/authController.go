package controllers

import (
	"dodolan/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// LoginCustomer godoc
// @Summary Login as customer.
// @Description Logging in to get jwt token to access customer api by roles.
// @Tags Auth
// @Param Body body LoginInput true "the body to login a customer"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /customer/login [post]
func LoginCustomer(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.Customer{}

	u.Username = input.Username
	u.Password = input.Password

	token, err := models.LoginCheckCustomer(u.Username, u.Password, db)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password customer is incorrect."})
		return
	}

	// user := map[string]string{
	// 	"username": u.Username,
	// 	"email":    u.Email,
	// }

	c.JSON(http.StatusOK, gin.H{"message": "login customer success", "username": u.Username, "token": token})

}

// LoginSeller godoc
// @Summary Login as seller.
// @Description Logging in to get jwt token to access seller api by roles.
// @Tags Auth
// @Param Body body LoginInput true "the body to login a seller"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /seller/login [post]
func LoginSeller(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.Seller{}

	u.Username = input.Username
	u.Password = input.Password

	token, err := models.LoginCheckSeller(u.Username, u.Password, db)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password seller is incorrect."})
		return
	}

	// user := map[string]string{
	// 	"username": u.Username,
	// 	"email":    u.Email,
	// }

	c.JSON(http.StatusOK, gin.H{"message": "login seller success", "username": u.Username, "token": token})

}
