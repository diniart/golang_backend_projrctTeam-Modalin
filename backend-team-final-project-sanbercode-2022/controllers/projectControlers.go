package controllers

import (
	"fintech/models"
	"fintech/utils/token"
	"strconv"
	"strings"

	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectsInput struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" gorm:"type:varchar(255)" validate:"omitempty,min=1,max=255"`
	Margin      int       `json:"margin"`
	Duration    int       `json:"duration"`
	Periode     int       `json:"periode"`
	Description string    `json:"description" gorm:"type:text"`
	Quantity    int       `json:"quantity"`
	Price       int       `json:"price"`
	Status      string    `json:"status"`
	DueDate     string    `json:"dueDate"`
	CategoryID  uuid.UUID `json:"category_id"`
	UserID      uuid.UUID `json:"user_id"`
}

type ProjectSold struct {
	ProjectID   uuid.UUID `gorm:"column:id"`
	Name        string    `gorm:"column:name"`
	Quantity    int       `gorm:"column:quantity"`
	Margin      int       `gorm:"column:margin"`
	Periode     int       `gorm:"column:periode"`
	Duration    int       `gorm:"column:duration"`
	Description string    `gorm:"column:description"`
	Price       int       `gorm:"column:price"`
	DueDate     time.Time `gorm:"column:due_date"`
	UserID      uuid.UUID `gorm:"column:user_id"`
	Sold        int       `gorm:"column:sold"`
	Status      string    `gorm:"column:status"`
}

type ProjectFiltered struct {
	ProjectID   uuid.UUID `gorm:"column:id"`
	Name        string    `gorm:"column:name"`
	Margin      int       `gorm:"column:margin"`
	Periode     int       `gorm:"column:periode"`
	Duration    int       `gorm:"column:duration"`
	Description string    `gorm:"column:description"`
	Price       int       `gorm:"column:price"`
	Status      string    `gorm:"column:status"`
	Sold        int       `gorm:"column:sold"`
	Quantity    int       `gorm:"column:quantity"`
	ImageURL    string    `gorm:"column:image_url"`
	CategoryID  uuid.UUID `gorm:"column:category_id"`
}

type ProjectCount struct {
	Status string `gorm:"column:status"`
	Count  int    `gorm:"column:count"`
}

// GetAllProjects godoc
// @Summary Get all projects.
// @Description Get a list of Projects.
// @Tags Projects
// @Produce json
// @Success 200 {object} []models.Projects
// @Router /projects [get]
func GetAllProjects(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var projects []models.Projects
	err := db.Preload("Category").Preload("Images").Preload("User.UserProfile").Find(&projects).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": projects})
}

// GetSomeProjects godoc
// @Summary Get some projects according to the limit.
// @Description Get some projects according to the limit. Add "?limit=" to the URL to specify the limit. Example: /projects/some?limit=4
// @Tags Projects
// @Produce json
// @Success 200 {object} []models.Projects
// @Router /projects/some [get]
func GetSomeProjects(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	lmt, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	var projects []models.Projects
	err := db.Limit(lmt).Preload("Category").Preload("Images").Preload("User.UserProfile").Find(&projects).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": projects})
}

// GetProjectsFiltered godoc
// @Summary Get projects according to the filter conditions.
// @Description Get some projects according to the filter conditions. Add "?limit=&status=" to the URL to specify the filter conditions. Example: /projects/filter?limit=4&status=running
// @Tags Projects
// @Produce json
// @Success 200 {object} []models.Projects
// @Router /projects/filter [get]
func GetProjectsFiltered(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	lmt, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))

	category := c.Request.URL.Query().Get("category")
	textCategory := []string{"%", category, "%"}
	joinedCategory := strings.Join(textCategory, "")

	status := c.Request.URL.Query().Get("status")
	textStatus := []string{"%", status, "%"}
	joinedStatus := strings.Join(textStatus, "")

	keyword := c.Request.URL.Query().Get("keyword")
	textKeyword := []string{"%", keyword, "%"}
	joinedKeyword := strings.Join(textKeyword, "")

	var projects []ProjectFiltered

	// err := db.Preload("Category").Preload("Images").Preload("User.UserProfile").Where("status like ?", joinedStatus).Where("category_id like ?", joinedCategory).Where("lower(name) like ?", joinedKeyword).Limit(lmt).Find(&projects).Error

	err := db.Raw(`select p2.id, p2.name, p2.margin, p2.duration , p2.periode , p2.description , p2.quantity ,p2.price,p2.status, y.image_url, x.sold, p2.category_id from (select w.id, sum(w.sold) as sold from 
	(select v.quantity as sold, p.id from 
		(select ci.quantity, ci.projects_id  from cart_items ci
			join shopping_carts sc on sc.id  = ci.shopping_cart_id 
			where sc.payment_status = 'success') v
	full join projects p on p.id = v.projects_id) w
group by w.id) x
full join  (select distinct on (i.projects_id) i.projects_id,i.image_url from images i order by i.projects_id  ) as y
on y.projects_id = x.id
join projects p2 on x.id= p2.id 
where status like ? and category_id like ? and lower(name) like ?
limit ?`, joinedStatus, joinedCategory, joinedKeyword, lmt).Scan(&projects).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": projects})
}

// CreateProjects godoc
// @Summary Create New Projects.
// @Description Creating a new projects.
// @Tags Projects
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param Body body ProjectsInput true "the body to create a new projects"
// @Produce json
// @Success 200 {object} models.Projects
// @Router /projects [post]
func CreateProjects(c *gin.Context) {
	// Validate input
	var input ProjectsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse time
	var layoutFormat = "2006-01-02 MST"
	value := fmt.Sprintf("%s WIB", input.DueDate)
	date, err := time.Parse(layoutFormat, value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get user id from token
	idUser, err := token.ExtractTokenID(c)
	if err != nil {
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

	// Create Projects
	projects := models.Projects{
		ID:          uuid.New(),
		Name:        input.Name,
		Margin:      input.Margin,
		Duration:    input.Duration,
		Periode:     input.Periode,
		Description: input.Description,
		Quantity:    input.Quantity,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
		UserID:      idUser,
		DueDate:     date,
	}

	db := c.MustGet("db").(*gorm.DB)
	err = db.Create(&projects).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	marginPerMonth := (float32(projects.Margin) + 5) / 12

	profitAmount := ((marginPerMonth * float32(projects.Periode)) * float32(projects.Quantity*projects.Price/100))

	installmentCount := projects.Duration / projects.Periode

	for i := 1; i <= installmentCount; i++ {

		installmentProfit := models.Installment{
			ID:         uuid.New(),
			Amount:     int(profitAmount),
			Status:     "Not-Paid",
			Type:       "Profit",
			ProjectsID: projects.ID,
		}

		// create installment only profit sharing margin
		err = db.Create(&installmentProfit).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	mainLoan := projects.Quantity * projects.Price

	mainInstallment := models.Installment{
		ID:         uuid.New(),
		Amount:     mainLoan,
		Status:     "Not-Paid",
		Type:       "Main",
		ProjectsID: projects.ID,
	}

	err = db.Create(&mainInstallment).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// create installment with main loan and last profit sharing margin

	c.JSON(http.StatusOK, gin.H{"data": projects})

}

// GetProjectsById godoc
// @Summary Get Projects.
// @Description Get an projects by id.
// @Tags Projects
// @Produce json
// @Param id path string true "Projects ID"
// @Success 200 {object} models.Projects
// @Router /projects/{id} [get]
func GetProjectsById(c *gin.Context) {
	var project models.Projects

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Set("gorm:auto_preload", true).Preload("Category").Preload("Images").Preload("User.UserProfile").Where("id = ?", c.Param("id")).First(&project).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project})

}

// GetProjectsByUserId godoc
// @Summary Get Projects.
// @Description Get an projects by user id.
// @Tags Projects
// @Produce json
// @Param id path string true "Projects ID"
// @Success 200 {object} models.Projects
// @Router /projects/investee/{userid} [get]
func GetProjectByUserId(c *gin.Context) {
	var project []models.Projects

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Set("gorm:auto_preload", true).Preload("Category").Preload("Images").Preload("User.UserProfile").Where("user_id = ?", c.Param("userid")).Find(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project})

}

// GetProjectsById godoc
// @Summary Get Projects.
// @Description Get an projects by id.
// @Tags Projects
// @Param Authorization header string true "Store Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Projects ID"
// @Success 200 {object} models.Projects
// @Router /projects/{id} [get]
func GetProjectsByUserId(c *gin.Context) {
	var projects []models.Projects

	// Extract Token ID User
	userID, _ := token.ExtractTokenID(c)

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Preload("Installment").Preload("Category").Preload("Images").Preload("User.UserProfile").Where("user_id = ?", userID).Find(&projects).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": projects})

}

// UpdateProjects godoc
// @Summary Update Projects.
// @Description Update Projects by ID.
// @Tags Projects
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param ID path string true "Projects ID"
// @Param Body body ProjectsInput true "the body to update projects"
// @Success 200 {object} models.Projects
// @Router /projects/{id} [patch]
func UpdateProjects(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var project models.Projects
	if err := db.Where("id = ?", c.Param("id")).First(&project).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// get role from token
	role, _ := token.ExtractTokenRole(c)
	// checking role
	if role != "investee" && role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}

	//  Validate input
	var input ProjectsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Parse time
	var layoutFormat = "2006-01-02 MST"
	value := fmt.Sprintf("%s WIB", input.DueDate)
	date, err := time.Parse(layoutFormat, value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var updatedInput models.Projects
	updatedInput.Name = input.Name
	updatedInput.Margin = input.Margin
	updatedInput.Duration = input.Duration
	updatedInput.Periode = input.Periode
	updatedInput.Description = input.Description
	updatedInput.Quantity = input.Quantity
	updatedInput.Price = input.Price
	updatedInput.CategoryID = input.CategoryID
	updatedInput.UserID = input.UserID
	updatedInput.DueDate = date
	updatedInput.Status = input.Status

	err = db.Model(&project).Updates(updatedInput).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": project})
}

// UpdateProjectsStatus godoc
// @Summary Update Projects Status.
// @Description Update Projects Status by ID.
// @Tags Projects
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param ID path string true "Projects ID"
// @Param Body body ProjectsInput true "the body to update projects"
// @Success 200 {object} models.Projects
// @Router /projects/status/{id} [patch]
func UpdateProjectsStatus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var project models.Projects
	if err := db.Where("id = ?", c.Param("id")).First(&project).Error; err != nil {
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
	var input ProjectsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.Projects
	updatedInput.Status = input.Status

	err := db.Model(&project).Updates(updatedInput).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project})
}

// DeleteProjects godoc
// @Summary Delete One Projects.
// @Description Delete a projects by id.
// @Tags Projects
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Projects ID"
// @Success 200 {object} map[string]boolean
// @Router /projects/{id} [delete]
func DeleteProjects(c *gin.Context) {
	//  Get Model if exist
	db := c.MustGet("db").(*gorm.DB)

	// get role from token
	role, _ := token.ExtractTokenRole(c)

	// checking role
	if role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission"})
		return
	}

	var project models.Projects
	if err := db.Where("id = ?", c.Param("id")).First(&project).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	db.Delete(&project)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// GetAllProjectSold godoc
// @Summary Get all projects sold.
// @Description Get all list of Projects that was sold.
// @Tags Projects
// @Produce json
// @Success 200 {object} []models.Projects
// @Router /projects/sold [get]
func GetAllProjectsSold(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)

	var projectSold []ProjectSold
	err := db.Raw(`select p.id, p.name, p.quantity,v.sold, p.margin, p.periode, p.duration, p.description, p.price, p.due_date, p.status, p.user_id  from (select ci.projects_id ,sum(ci.quantity) as sold from cart_items ci 
	left join shopping_carts sc on ci.shopping_cart_id  = sc.id 
	where sc.payment_status = 'success'
	group by ci.projects_id) v full join projects p on p.id = v.projects_id
order by p.created_at desc`).Scan(&projectSold).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": projectSold})
}

// GetAllProjectSoldByParamID godoc
// @Summary Get projects sold by param id.
// @Description Get quantity Projects that was sold by using param id.
// @Tags Projects
// @Produce json
// @Success 200 {object} []models.Projects
// @Router /projects/sold/:id [get]
func GetProjectsSoldByParamID(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)

	var projectSold []ProjectSold
	err := db.Raw(`select v.*, p.name, p.quantity, p.margin, p.periode, p.duration, p.description, p.price, p.due_date, p.status, p.user_id  from (select ci.projects_id ,sum(ci.quantity) as sold from cart_items ci 
	left join shopping_carts sc on ci.shopping_cart_id  = sc.id 
	where sc.payment_status = 'success'
	group by ci.projects_id) v 
	join projects p on p.id = v.projects_id
	where p.id = ?`, c.Param("id")).Scan(&projectSold).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": projectSold})
}

// GetAllProjectBuyByID godoc
// @Summary Get all projects buy by token investor.
// @Description Get quantity Projects buy that was bought by using token .
// @Tags Projects
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []models.Projects
// @Router /projects/buy [get]
func GetProjectsBuyByToken(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)

	// Extract Token ID User
	userID, _ := token.ExtractTokenID(c)

	var projectSold []ProjectSold
	err := db.Raw(`select p.id, v.sold, p.name, p.quantity, p.margin, p.periode, p.duration, p.description, p.price, p.due_date, p.status, p.user_id  from (select projects_id, user_id as investor, sum(quantity) as sold from cart_items ci 
	left join shopping_carts sc on ci.shopping_cart_id  = sc.id 
	where payment_status = 'success'
	group by projects_id, user_id) v
	join projects p on p.id = v.projects_id
	where investor = ?`, userID).Scan(&projectSold).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": projectSold})
}

// GetProjectsCount godoc
// @Summary Get projects count group by project status.
// @Description Get Projects count group by project status .
// @Tags Projects
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []ProjectCount
// @Router /projects/count [get]
func GetProjectsCount(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)

	// Extract Token ID User
	userID, _ := token.ExtractTokenID(c)

	var projectCount []ProjectCount
	err := db.Raw(`select count(status), p.status from projects p where user_id = ?	group by status`, userID).Scan(&projectCount).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": projectCount})
}
