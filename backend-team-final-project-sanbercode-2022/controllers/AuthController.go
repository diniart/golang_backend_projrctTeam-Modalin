package controllers

import (
	"fintech/models"
	"fintech/utils/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePassInput struct {
	Password    string `json:"password" binding:"required"`
	PasswordNew string `json:"password_new" binding:"required"`
}

// LoginUser godoc
// @Summary Login as as user.
// @Description Logging in to get jwt token to access admin or user api by roles.
// @Tags Authentication
// @Param Body body LoginInput true "the body to login a user"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /login [post]
func Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := models.User{}

	u.Email = input.Email
	u.Password = input.Password

	token, role, status, strID, err := models.LoginCheck(u.Email, u.Password, db)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or password is incorrect."})
		return
	}

	if status != "aktif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account deactived"})
		return
	}

	user := map[string]string{
		"id":    strID,
		"email": u.Email,
		"role":  role,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login succes", "user": user, "status": status, "token": token})

}

// Register godoc
// @Summary Register a user.
// @Description registering a user from public access.
// @Tags Authentication
// @Param Body body RegisterInput true "the body to register a user"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /register [post]
func Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}
	u.Email = input.Email
	u.Role = input.Role
	u.Password = input.Password

	_, err := u.SaveUser(db)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := map[string]string{
		"email": input.Email,
		"role":  input.Role,
	}

	c.JSON(http.StatusOK, gin.H{"message": " regristation success", "user": user})
}

// ChangePassword godoc
// @Summary Update Password.
// @Description Update Password by id.
// @Tags Password
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "Password id"
// @Success 200 {object} models.User
// @Router /userProfile/{id} [patch]
func ChangePassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input ChangePassInput

	// validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Extract Token ID User
	userID, _ := token.ExtractTokenID(c)

	u := models.User{}
	u.ID = userID
	u.Password = input.Password
	// Check Password
	err := models.PasswordCheck(u.ID.String(), u.Password, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is incorrect."})
		return
	}

	// Get model if exist
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	// hash the password

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.PasswordNew), bcrypt.DefaultCost)

	var updatedInput models.User
	updatedInput.Password = string(hashedPassword)
	updatedInput.UpdatedAt = time.Now()

	err = db.Model(&user).Updates(updatedInput).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error updates data."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
