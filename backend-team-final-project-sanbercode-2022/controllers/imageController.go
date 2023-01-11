package controllers

import (
	"fintech/models"
	"fintech/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type imageInput struct {
	ImageURL   string    `json:"images_url"`
	ProjectsID uuid.UUID `json:"projects_id"`
}

// CreateNewImage godoc
// @Summary Create New Image.
// @Description Admin can create a new image.
// @Tags Images
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Accept json
// @Produce json
// @Success 200 {object} models.Images
// @Router /images [post]
func CreateNewImage(c *gin.Context) {

	// Validate image input
	var input imageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get role from token
	role, _ := token.ExtractTokenRole(c)
	// checking role
	if role != "investee" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}

	// Create Image
	image := models.Images{
		ID:         uuid.New(),
		ImageURL:   input.ImageURL,
		ProjectsID: input.ProjectsID,
	}

	db := c.MustGet("db").(*gorm.DB)
	err := db.Create(&image).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": image})
}

// GetImagesByProjectId godoc
// @Summary Get Images by Project ID.
// @Description Get Images by Project id.
// @Tags Images
// @Produce json
// @Param id path string true "Image id"
// @Success 200 {object} models.Images
// @Router /images/{id} [get]
func GetImagesByProjectId(c *gin.Context) { // Get model if exist

	var image models.Images

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("projects_id = ?", c.Param("id")).Find(&image).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": image})
}

// UpdateImage godoc
// @Summary Update Image.
// @Description Update Image by id.
// @Tags Images
// @Produce json
// @Param id path string true "Image id"
// @Param Body body imageInput true "the body to update images"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} models.Images
// @Router /images/{id} [patch]
func UpdateImage(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model
	var image models.Images
	if err := db.Where("id = ?", c.Param("id")).First(&image).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// get role from token
	role, _ := token.ExtractTokenRole(c)
	// checking role
	if role != "investee" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}

	// Validate input
	var input imageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.Images
	updatedInput.ImageURL = input.ImageURL

	err := db.Model(&image).Updates(updatedInput).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": image})
}

// DeleteImage godoc
// @Summary Delete one Image.
// @Description Delete a Image by id.
// @Tags Images
// @Produce json
// @Param id path string true "Image id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]boolean
// @Router /images/{id} [delete]
func DeleteImage(c *gin.Context) {
	// Get model
	db := c.MustGet("db").(*gorm.DB)
	var image models.Images
	if err := db.Where("id = ?", c.Param("id")).First(&image).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// get role from token
	role, _ := token.ExtractTokenRole(c)
	// checking role
	if role != "investee" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}

	err := db.Delete(&image).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "deleted"})
}

// function not used :

// // GetAllImages godoc
// // @Summary Get all Images.
// // @Description Get a list of Images.
// // @Tags Images
// // @Produce json
// // @Success 200 {object} []models.Images
// // @Router /images [get]
// func GetAllImages(c *gin.Context) {
// 	// get db from gin context
// 	db := c.MustGet("db").(*gorm.DB)
// 	var images []models.Images
// 	db.Find(&images)

// 	c.JSON(http.StatusOK, gin.H{"data": images})
// }
