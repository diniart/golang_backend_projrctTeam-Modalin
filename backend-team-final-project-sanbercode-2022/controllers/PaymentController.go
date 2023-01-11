package controllers

import (
	"fintech/models"
	"fintech/utils/token"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type PaymentInput struct {
	OrderID  string `json:"order_id"`
	GrossAmt int    `json:"gross_amt"`
}

type NotificationInput struct {
	TransactionID     string `json:"transaction_id"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
}

var s snap.Client

func CreateMidtransTransaction(c *gin.Context) {
	s.New("SB-Mid-server-L3OGC3dH1nFYqq5jbBpfymQB", midtrans.Sandbox)
	// Validate input
	var input PaymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract Token ID User
	userID, _ := token.ExtractTokenID(c)

	// get user data
	var user models.User
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	// order id bisa diisi shopping cart id or project id
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  input.OrderID,
			GrossAmt: int64(input.GrossAmt),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: user.Email,
			FName: user.Role,
		},
	}

	// // Send request to Midtrans Snap API
	resp, err := s.CreateTransaction(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.GetMessage()})
		return
	}

	// // find models
	// var cart models.ShoppingCart
	// if errs := db.Where("id = ?", input.OrderID).First(&cart).Error; errs != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"token": resp.Token, "redirect_url": resp.RedirectURL})

}

func PaymentNotification(c *gin.Context) {
	var d = coreapi.Client{}
	d.New("SB-Mid-server-L3OGC3dH1nFYqq5jbBpfymQB", midtrans.Sandbox)

	db := c.MustGet("db").(*gorm.DB)

	var input NotificationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// order id bisa berisi shopping cart id atau cicilan id
	orderID, _ := uuid.Parse(input.OrderID)
	// parsing data
	transactionID, _ := uuid.Parse(input.TransactionID)
	adminID, _ := uuid.Parse("502aae43-40f4-4367-a973-1e47a5d18957")
	amount, _ := strconv.ParseFloat(input.GrossAmount, 32)

	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, e := d.CheckTransaction(input.OrderID)
	if e != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found!"})
		return
	} else {

		if transactionStatusResp != nil {
			// pengecekan apakah order id itu dari shopping cart id atau dari project id

			var cart models.ShoppingCart
			var installment models.Installment

			errs := db.Where("id = ?", orderID).First(&cart).Error

			// kalau gak ketemu di cart maka cek di installment
			if errs != nil {
				errs = db.Where("id = ?", orderID).First(&installment).Error
				if errs != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
					return
				}

				// updates status in installment data to paid
				var updatedStatusInstallment models.Installment

				var projects models.Projects
				errs = db.Where("id = ?", installment.ProjectsID).First(&projects).Error
				if errs != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
					return
				}

				// 5. Do set transaction status based on response from check transaction status
				if transactionStatusResp.TransactionStatus == "settlement" {
					// TODO set transaction status on your databaase to 'success'
					updatedStatusInstallment.Status = "Paid"

					// create transaction

					transactionData := models.Transaction{
						ID:            transactionID,
						UserID:        adminID,
						Credit:        int(amount),
						Sender:        projects.UserID,
						Status:        "Cicilan",
						ProjectID:     installment.ProjectsID,
						InstallmentID: orderID,
					}

					err := db.Create(&transactionData).Error
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}

				} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
					// TODO set transaction status on your databaase to 'failure'
					updatedStatusInstallment.Status = "failure"

				} else if transactionStatusResp.TransactionStatus == "pending" {
					// TODO set transaction status on your databaase to 'pending' / waiting payment
					updatedStatusInstallment.Status = "pending"

				} else {
					updatedStatusInstallment.Status = "deny"
				}

				errs = db.Model(&installment).Updates(updatedStatusInstallment).Error
				if errs != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": errs.Error()})
					return
				}

			} else {
				var updatedInput models.ShoppingCart

				// 5. Do set transaction status based on response from check transaction status
				if transactionStatusResp.TransactionStatus == "settlement" {
					// TODO set transaction status on your databaase to 'success'
					updatedInput.PaymentStatus = "success"
					// create transaction

					transactionData := models.Transaction{
						ID:             transactionID,
						UserID:         adminID,
						Credit:         int(amount),
						Sender:         cart.UserID,
						Status:         "Pembayaran",
						ShoppingCartID: cart.ID,
					}

					err := db.Create(&transactionData).Error
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}

				} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
					// TODO set transaction status on your databaase to 'failure'
					updatedInput.PaymentStatus = "failure"

				} else if transactionStatusResp.TransactionStatus == "pending" {
					// TODO set transaction status on your databaase to 'pending' / waiting payment
					updatedInput.PaymentStatus = "pending"

				} else {
					updatedInput.PaymentStatus = "deny"

				}

				errs = db.Model(&cart).Updates(updatedInput).Error
				if errs != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": errs.Error()})
					return
				}

			}

		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})

}
