package controllers

import (
	"fintech/models"
	"fintech/utils/token"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cartItemsInput struct {
	Quantity       int       `json:"quantity"`
	ProjectsID     uuid.UUID `json:"projectID"`
	ShoppingCartID uuid.UUID `json:"shoppingCartID"`
}

// GetAllCartItems godoc
// @Summary Get all cartItem.
// @Description Get a list of CartItems.
// @Tags CartItems
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []models.CartItems
// @Router /cartItems [get]
func GetAllCartItems(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var cartItems []models.CartItems
	err := db.Find(&cartItems).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cartItems})

}

// GetCartItemsById godoc
// @Summary Get CartItems.
// @Description Get an CartItems by id.
// @Tags CartItems
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "CartItems ID"
// @Success 200 {object} models.CartItems
// @Router /cartItems/{id} [get]
func GetCartItemsById(c *gin.Context) {
	var cartItems models.CartItems

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Preload("Project").Preload("ShoppingCart").Where("id = ?", c.Param("id")).First(&cartItems).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cartItems})

}

// // CreateCartItems godoc
// // @Summary Create New CartItems.
// // @Description Creating a new CartItems.
// // @Tags CartItems
// // @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// // @Security BearerToken
// // @Produce json
// // @Success 200 {object} models.CartItems
// // @Router /cartItems [post]
// func CreateCartItems(c *gin.Context) {
// 	var input cartItemsInput

// 	if err := c.ShouldBind(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// get role from token
// 	role, _ := token.ExtractTokenRole(c)
// 	// checking role
// 	if role != "investor" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
// 		return
// 	}

	
// 	cartItems := models.CartItems{ID: uuid.New(), Quantity: input.Quantity, ProjectsID: input.ProjectsID, ShoppingCartID: input.ShoppingCartID}
// 	db := c.MustGet("db").(*gorm.DB)
// 	err := db.Create(&cartItems).Error
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": cartItems})

// }

// CreateCartItems godoc
// @Summary Create New CartItems.
// @Description Creating a new CartItems.
// @Tags CartItems
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} models.CartItems
// @Router /cartItems [post]
func CreateCartItems(c *gin.Context) {
	var input cartItemsInput
	db := c.MustGet("db").(*gorm.DB)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	// get role from token
	role, _ := token.ExtractTokenRole(c)
	// checking role
	if role != "investor" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}
	// check ada data project sama tidak
	var cartItem models.CartItems
	erro := db.Where("projects_id = ?",input.ProjectsID ).Where("shopping_cart_id = ?", input.ShoppingCartID).First(&cartItem).Error; 
	if erro != nil {
		cartItems := models.CartItems{ID: uuid.New(), Quantity: input.Quantity, ProjectsID: input.ProjectsID, ShoppingCartID: input.ShoppingCartID}

		err := db.Create(&cartItems).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": cartItems})

	}
	
	fmt.Println("qty>>>",cartItem.Quantity)
	var updatedInput models.CartItems
	updatedInput.Quantity = cartItem.Quantity + input.Quantity

	err := db.Model(&cartItem).Updates(updatedInput).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cartItem})

	


}

// UpdateCartItems godoc
// @Summary Update CartItems.
// @Description Update CartItems by id.
// @Tags CartItems
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "CartItems id"
// @Success 200 {object} models.CartItems
// @Router /CartItems/{id} [patch]
func UpdateCartItems(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	// Extract Token ID User
	userID, _ := token.ExtractTokenID(c)

	// get role from token
	role, _ := token.ExtractTokenRole(c)
	// checking role
	if role != "investor" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}

	// Get model if exist
	var cartItem models.CartItems
	if err := db.Where("id = ?", userID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input cartItemsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.CartItems
	updatedInput.Quantity = input.Quantity
	updatedInput.ProjectsID = input.ProjectsID
	updatedInput.ShoppingCartID = input.ShoppingCartID

	// updatedInput.UpdatedAt = time.Now()

	err := db.Model(&cartItem).Updates(updatedInput).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cartItem})
}

// DeleteCartItem godoc
// @Summary Delete one Cart item.
// @Description Delete a Cart item by id.
// @Tags CartItems
// @Produce json
// @Param id path string true "CartItem id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]boolean
// @Router /cartItems/{id} [delete]
func DeleteCartItem(c *gin.Context) {
	// Get model
	db := c.MustGet("db").(*gorm.DB)
	var cartItem models.CartItems
	if err := db.Where("id = ?", c.Param("id")).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&cartItem)

	c.JSON(http.StatusOK, gin.H{"data": true})
}