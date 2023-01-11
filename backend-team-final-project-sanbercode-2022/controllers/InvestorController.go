package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Investor struct {
	UserID        uuid.UUID `gorm:"column:user_id"`
	ProjectID     uuid.UUID `gorm:"column:projects_id"`
	Email         string    `gorm:"column:email"`
	ProjectName   string    `gorm:"column:name"`
	Margin        int       `gorm:"column:margin"`
	Duration      int       `gorm:"column:duration"`
	Periode       int       `gorm:"column:periode"`
	Quantity      int       `gorm:"column:quantity"`
	Price         int       `gorm:"column:price"`
	TotalQuantity int       `gorm:"column:total"`
}

// GetInvestorDataByProjectID godoc
// @Summary Get All Investor included in a project By ProjectID
// @Description Get all investor included in a project by project id.
// @Tags Investor
// @Param Authorization header string true "Store Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Projects ID"
// @Success 200 {object} models.Installment
// @Router /investor/{id} [get]
func GetAllInvestorByProjectID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var investor []Investor
	err := db.Raw(`select v.user_id, v.projects_id, u.email, p."name", p.margin, p.duration, p.quantity as total, p.periode, v.quantity, p.price  from (select sc.user_id, sum(quantity) as quantity,ci.projects_id  from cart_items ci 
	left join shopping_carts sc on ci.shopping_cart_id  = sc.id
	where ci.projects_id = ? and sc.payment_status = ?
	group by sc.user_id, ci.projects_id ) v
	left join users u on v.user_id  = u.id 
	left join projects p on v.projects_id  = p.id `, c.Param("id"), "success").Scan(&investor).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": investor})
}
