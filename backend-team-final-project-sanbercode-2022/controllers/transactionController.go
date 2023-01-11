package controllers

import (
	"fintech/models"
	"fintech/utils/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateNewTransaction godoc
// @Summary Create New Transaction Data By Role Investor and Investee.
// @Description Investor and Investee can create a new transaction data.
// @Tags Transaction
// @Param Authorization header string true "Store Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Accept json
// @Produce json
// @Success 200 {object} models.Transaction
// @Router /transaction [post]
func CreateTransaction(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)

	// Validate transaction input
	var input models.TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract Token ID User
	userID, _ := token.ExtractTokenID(c)

	// Create Transaction
	transaction := models.Transaction{
		ID:        uuid.New(),
		Debit:     input.Debit,
		Credit:    input.Credit,
		Sender:    userID,
		UserID:    input.Sender,
		Status:    input.Status,
		ProjectID: input.ProjectID,
	}

	// Save Transaction to database
	err := db.Create(&transaction).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Return Data Transaction
	c.JSON(http.StatusOK, gin.H{"data": "transaction success"})
}

// GetAllTransactionByUserID godoc
// @Summary investor and investee can get all transaction history data by user id in token.
// @Description Investor and Investee can get all transaction data by user ID
// @Tags Transaction
// @Param Authorization header string true "Store Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} models.Transaction
// @Router /transaction [get]
func GetAllTransactionByUserID(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)

	// Extract Token ID User
	userID, _ := token.ExtractTokenID(c)

	// Get All Transaction Data By User ID
	var transactions []models.Transaction
	err := db.Where("user_id = ?", userID).Order("created_at desc").Find(&transactions).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return All Data Transaction
	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

// GetTransactionsFiltered godoc
// @Summary Get transactions according to the filter conditions.
// @Description Get some projects according to the filter conditions. Add "?status=" to the URL to specify the filter conditions. Example: /transaction/filter?status=ProsesPencairan
// @Tags Transaction
// @Param Authorization header string true "Store Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []models.Transaction
// @Router /transaction/filter [get]
func GetTransactionsFiltered(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)

	// get role from token
	role, _ := token.ExtractTokenRole(c)
	// checking role
	if role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}

	status := c.Request.URL.Query().Get("status")
	textStatus := []string{"%", status, "%"}
	joinedStatus := strings.Join(textStatus, "")

	var transactions []models.Transaction

	err := db.Preload("User.UserProfile").Where("status like ?", joinedStatus).Find(&transactions).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

// UpdateTransaction godoc
// @Summary Update Transactions.
// @Description Update Transactions by ID.
// @Tags Transaction
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param ID path string true "Transaction ID"
// @Success 200 {object} models.Transaction
// @Router /transaction/{id} [patch]
func UpdateTransactions(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var transaction models.Transaction
	if err := db.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
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

	//  Validate input
	var input models.TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.Transaction
	updatedInput.Status = input.Status

	err := db.Model(&transaction).Updates(updatedInput).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": transaction})
}
