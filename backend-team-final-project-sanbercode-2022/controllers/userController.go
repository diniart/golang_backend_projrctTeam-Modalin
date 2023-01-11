package controllers

import (
	"fintech/models"
	"fintech/utils/token"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserInput struct {
	Email    string `json:"email" `
	Role     string `json:"role" `
	Password string `json:"password" `
	Status   string `json:"status" `
}

// UpdateUser godoc
// @Summary Update User.
// @Description Update User by id.
// @Tags admin
// @Produce json
// @Param id path string true "User id"
// @Param Body body UserInput true "the body to update user"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} models.User
// @Router /admin/user/{id} [patch]
func UpdateUser(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// get role from token
	role, _ := token.ExtractTokenRole(c)
	// checking role
	if role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}

	// Validate input
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash the password

	// hashedPassword,_ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	var updatedInput models.User
	// updatedInput.Email = input.Email
	// updatedInput.Password = string(hashedPassword)
	updatedInput.Status = input.Status
	// updatedInput.UpdatedAt = time.Now()

	err := db.Model(&user).Updates(updatedInput).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
