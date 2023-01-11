package controllers

import (
	"fintech/models"
	"fintech/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllInstallmentByProjectID godoc
// @Summary Get All Installment By ProjectID
// @Description Get all installment by project id.
// @Tags Installment
// @Param Authorization header string true "Store Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Projects ID"
// @Success 200 {object} models.Installment
// @Router /installment/{id} [get]
func GetAllInstallmentByProjectID(c *gin.Context) {
	var installment []models.Installment

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("projects_id = ?", c.Param("id")).Find(&installment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": installment})

}

// UpdateInstallmentStatus godoc
// @Summary Update Installment Status.
// @Description Update Installment Status by Installment ID.
// @Tags Installment
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param ID path string true "Installment ID"
// @Param Body body models.InstallmentInput true "the body to update installment status"
// @Success 200 {object} models.Installment
// @Router /installment/status/{id} [patch]
func UpdateInstallmentStatus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get model if exist
	var installment models.Installment
	if err := db.Where("id = ?", c.Param("id")).First(&installment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	// get role from token
	role, _ := token.ExtractTokenRole(c)
	// checking role
	if role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}

	//  Validate input
	var input models.InstallmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.Installment
	updatedInput.Status = input.Status

	err := db.Model(&installment).Updates(updatedInput).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": installment})
}
