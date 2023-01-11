package controllers

import (
	"fintech/models"
	"fintech/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cartInput struct {
	PaymentStatus string `json:"payment_status"`
	// CartItems
}

// GetAllCarts godoc
// @Summary Get all Carts.
// @Description Get a list of Carts.
// @Tags Cart
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []models.ShoppingCart
// @Router /carts [get]
func GetAllCarts(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var carts []models.ShoppingCart
	err := db.Find(&carts).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": carts})
}

// CreateNewCart godoc
// @Summary Create New Cart.
// @Description Admin can create a new cart.
// @Tags Cart
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Accept json
// @Produce json
// @Success 200 {object} models.ShoppingCart
// @Router /cart [post]
func CreateNewCart(c *gin.Context) {

	// Validate cart input
	var input cartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	UserID, _ := token.ExtractTokenID(c)
	// Create Cart
	cart := models.ShoppingCart{
		ID:            uuid.New(),
		PaymentStatus: input.PaymentStatus,
		UserID:        UserID,
	}
	// get role from token
	role, _ := token.ExtractTokenRole(c)

	// checking role
	if ( role != "investor" ) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	err := db.Create(&cart).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// GetCartByUserId godoc
// @Summary Get Cart.
// @Description Get an Cart by User id.
// @Tags Cart
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} models.ShoppingCart
// @Router /cart [get]
func GetCartByUserId(c *gin.Context) {
	var cart []models.ShoppingCart

	UserID, _ := token.ExtractTokenID(c)

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Preload("CartItems.Projects.Images").Where("user_id = ?", UserID).Find(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// GetCartOrder godoc
// @Summary Get Cart Payment Status "order".
// @Description Get Cart Payment Status "order".
// @Tags Cart
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} models.ShoppingCart
// @Router /cart-order [get]
func GetCartOrder(c *gin.Context) {
	var cart []models.ShoppingCart

	UserID, _ := token.ExtractTokenID(c)

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Preload("CartItems.Projects.Images").Where("user_id = ? AND payment_status =?", UserID, "order").Find(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// UpdateCart godoc
// @Summary Update Cart.
// @Description Update Cart by id.
// @Tags Cart
// @Produce json
// @Param id path string true "Cart id"
// @Param Body body cartInput true "the body to update cart"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} models.ShoppingCart
// @Router /cart/{id} [patch]
func UpdateCart(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model
	var cart models.ShoppingCart
	if err := db.Where("id = ?", c.Param("id")).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input cartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get role from token
	role, _ := token.ExtractTokenRole(c)

	// checking role
	if ( role != "investor" && role != "admin" ) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}

	var updatedInput models.ShoppingCart
	updatedInput.PaymentStatus = input.PaymentStatus

	err := db.Model(&cart).Updates(updatedInput).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// // DeleteCart godoc
// // @Summary Delete one Cart.
// // @Description Delete a Cart by id.
// // @Tags Cart
// // @Produce json
// // @Param id path string true "Cart id"
// // @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// // @Security BearerToken
// // @Success 200 {object} map[string]boolean
// // @Router /cart/{id} [delete]
// func DeleteCart(c *gin.Context) {
// 	// Get model
// 	db := c.MustGet("db").(*gorm.DB)
// 	var cart models.ShoppingCart
// 	if err := db.Where("id = ?", c.Param("id")).First(&cart).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
// 		return
// 	}

// 	db.Delete(&cart)

// 	c.JSON(http.StatusOK, gin.H{"data": true})
// }
