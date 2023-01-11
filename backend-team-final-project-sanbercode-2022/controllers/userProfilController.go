package controllers

import (
	"fintech/models"
	"fintech/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userProfileInput struct {
	Fullname          string `json:"fullName"`
	Phone             string `json:"phone"`
	KTP               int    `json:"kTP"`
	Address           string `json:"address"`
	City              string `json:"city"`
	Province          string `json:"province"`
	Gender            string `json:"gender"`
	ProfileURL        string `json:"profile_url"`
	BankAccountNumber int    `json:"bank_account_number"`
	BankName          string `json:"bank_name"`
}

// GetAllUserProfile godoc
// @summary Get a list of UserProfile.
// @Description Get a list of UserProfile.
// @Tags admin
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []models.UserProfile
// @Router /admin/users [get]
func GetAllUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var users []models.User
	err := db.Preload("UserProfile").Find(&users).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetAllUserProfile by Investor godoc
// @summary Get a list of UserProfile by Investor.
// @Description Get a list of UserProfile by Investor.
// @Tags admin
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []models.UserProfile
// @Router /admin/users/investor [get]
func GetAllUserInvestor(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var users []models.User
	err := db.Preload("UserProfile").Where("role = ?", "investor").Find(&users).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetAllUserProfile by Investee godoc
// @summary Get a list of UserProfile by Investee.
// @Description Get a list of UserProfile by Investee.
// @Tags admin
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []models.UserProfile
// @Router /admin/users/investee [get]
func GetAllUserInvestee(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var users []models.User
	err := db.Preload("UserProfile").Where("role = ?", "investee").Find(&users).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetUserProfilById godoc
// @Summary Get UserProfile.
// @Description Get an UserProfile by id.
// @Tags UserProfile
// @Produce json
// @Param id path string true "UserProfile id"
// @Success 200 {object} models.UserProfile
// @Router /userProfile/{id} [get]
func GetUserProfileByParamId(c *gin.Context) { // Get model if exist
	var user models.UserProfile

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// GetUserProfilById godoc
// @Summary Get UserProfile.
// @Description Get an UserProfile by id.
// @Tags UserProfile
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "UserProfile id"
// @Success 200 {object} models.UserProfile
// @Router /userProfile/{id} [get]
func GetUserProfileById(c *gin.Context) { // Get model if exist
	var user models.UserProfile

	// Extract Token ID User
	userID, _ := token.ExtractTokenID(c)

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUserProfile godoc
// @Summary Update UserProfile.
// @Description Update UserProfile by id.
// @Tags UserProfile
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "UserProfile id"
// @Success 200 {object} models.UserProfile
// @Router /userProfile/{id} [patch]
func UpdateUserProfile(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	// Extract Token ID User
	userID, _ := token.ExtractTokenID(c)

	// Get model if exist
	var user models.UserProfile
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input userProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.UserProfile
	updatedInput.Fullname = input.Fullname
	updatedInput.Phone = input.Phone
	updatedInput.KTP = input.KTP
	updatedInput.Address = input.Address
	updatedInput.City = input.City
	updatedInput.Province = input.Province
	updatedInput.Gender = input.Gender
	updatedInput.ProfileURL = input.ProfileURL
	updatedInput.BankAccountNumber = input.BankAccountNumber
	updatedInput.BankName = input.BankName

	err := db.Model(&user).Updates(updatedInput).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// Function / Feature Not Used :

// // DeleteUserProfile godoc
// // @Summary Delete one UserProfile.
// // @Description Delete a UserProfile by id.
// // @Tags UserProfile
// // @Produce json
// // @Param id path string true "UserProfile id"
// // @Success 200 {object} map[string]boolean
// // @Router /userProfile/{id} [delete]
// func DeleteUserProfile(c *gin.Context) {
// 	// Get model if exist
// 	db := c.MustGet("db").(*gorm.DB)
// 	var user models.UserProfile
// 	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
// 		return
// 	}

// 	db.Delete(&user)

// 	c.JSON(http.StatusOK, gin.H{"data": true})
// }

// // CreateUserProfile godoc
// // @Summary Create New UserProfile.
// // @Description Creating a new UserProfile.
// // @Tags UserProfile
// // @Security BearerToken
// // @Produce json
// // @Success 200 {object} models.UserProfile
// // @Router /userProfile [post]
// func CreateUserProfile(c *gin.Context) {
// 	var input userProfileInput
// 	if err := c.ShouldBind(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	userProfil := models.UserProfile{Fullname: input.Fullname, Phone: input.Phone, KTP: input.KTP, Address: input.Address, City: input.City, Province: input.Province, Gender: input.Gender, ProfileURL: input.ProfileURL}
// 	db := c.MustGet("db").(*gorm.DB)
// 	db.Create(&userProfil)
// 	c.JSON(http.StatusOK, gin.H{"data": userProfil})

// }
