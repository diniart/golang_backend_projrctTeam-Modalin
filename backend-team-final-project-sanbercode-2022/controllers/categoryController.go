package controllers

import (
	"fintech/models"
	"fintech/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryInput struct {
	Category string `json:"category"`
}

// GetAllCategories godoc
// @Summary Get all Categories.
// @Description Get a list of Category.
// @Tags Category
// @Produce json
// @Success 200 {object} []models.Category
// @Router /category [get]
func GetAllCategories(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var categories []models.Category
	err := db.Find(&categories).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// CreateNewCategory godoc
// @Summary Create New Category.
// @Description Admin can create a new category.
// @Tags Category
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Accept json
// @Produce json
// @Success 200 {object} models.Category
// @Router /category [post]
func CreateNewCategory(c *gin.Context) {

	// Validate category input
	var input categoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// get role from token
	role, _ := token.ExtractTokenRole(c)
	// checking role
	if role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}

	// Create Category
	category := models.Category{
		ID:       uuid.New(),
		Category: input.Category,
	}

	db := c.MustGet("db").(*gorm.DB)
	err := db.Create(&category).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// GetCategoryById godoc
// @Summary Get Category.
// @Description Get a Category by id.
// @Tags Category
// @Produce json
// @Param id path string true "Category id"
// @Success 200 {object} models.Category
// @Router /category/{id} [get]
func GetCategoryById(c *gin.Context) { // Get model if exist
	var category models.Category

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// UpdateCategory godoc
// @Summary Update Category.
// @Description Update Category by id.
// @Tags Category
// @Produce json
// @Param id path string true "Category id"
// @Param Body body categoryInput true "the body to update category"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} models.Category
// @Router /category/{id} [patch]
func UpdateCategory(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model
	var category models.Category
	if err := db.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
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
	var input categoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.Category
	updatedInput.Category = input.Category

	err := db.Model(&category).Updates(updatedInput).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// DeleteCategory godoc
// @Summary Delete one Category.
// @Description Delete a Category by id.
// @Tags Category
// @Produce json
// @Param id path string true "Category id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]boolean
// @Router /category/{id} [delete]
func DeleteCategory(c *gin.Context) {
	// Get model
	db := c.MustGet("db").(*gorm.DB)
	var category models.Category
	if err := db.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
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

	err := db.Delete(&category).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
