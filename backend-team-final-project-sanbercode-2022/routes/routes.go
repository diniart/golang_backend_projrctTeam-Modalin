package routes

import (
	"fintech/controllers"
	"fintech/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{"localhost"})

	// CORS Setting
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"}
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true
	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS")
	r.Use(cors.New(corsConfig))

	// set db to gin context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	// JWT Setup
	jwt := r.Group("/")
	jwt.Use(middlewares.JwtAuthMiddleware())

	// ROUTES SETUP
	// AUTH
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	jwt.POST("/change_password", controllers.ChangePassword)

	//Category Routes
	r.GET("/category", controllers.GetAllCategories)
	r.GET("/category/:id", controllers.GetCategoryById)

	jwt.POST("/category", controllers.CreateNewCategory)
	jwt.PATCH("/category/:id", controllers.UpdateCategory)
	jwt.DELETE("/category/:id", controllers.DeleteCategory)

	//Images Routes
	r.GET("/images/:id", controllers.GetImagesByProjectId)

	jwt.POST("/images", controllers.CreateNewImage)
	jwt.PATCH("/images/:id", controllers.UpdateImage)
	jwt.DELETE("/images/:id", controllers.DeleteImage)

	//Shopping Cart Routes
	jwt.GET("/carts", controllers.GetAllCarts)
	jwt.GET("/cart", controllers.GetCartByUserId)
	jwt.POST("/cart", controllers.CreateNewCart)
	jwt.PATCH("/cart/:id", controllers.UpdateCart)
	jwt.GET("/cart-order", controllers.GetCartOrder)

	// User Profile Routes
	r.GET("/userProfile/:id", controllers.GetUserProfileByParamId)
	jwt.GET("/userProfile", controllers.GetUserProfileById)

	jwt.PATCH("/userProfile", controllers.UpdateUserProfile)

	// Transaction Routes
	jwt.GET("/transaction", controllers.GetAllTransactionByUserID)
	jwt.POST("/transaction", controllers.CreateTransaction)
	jwt.GET("/transaction/filter", controllers.GetTransactionsFiltered)
	jwt.PATCH("/transaction/:id", controllers.UpdateTransactions)

	// Admin Routes
	jwt.GET("/admin/users", controllers.GetAllUser)
	jwt.GET("/admin/users/investor", controllers.GetAllUserInvestor)
	jwt.GET("/admin/users/investee", controllers.GetAllUserInvestee)
	jwt.PATCH("/admin/user/:id", controllers.UpdateUser)

	jwt.GET("/investor/:id", controllers.GetAllInvestorByProjectID)

	// PROJECT
	r.GET("/projects", controllers.GetAllProjects)
	r.GET("/projects/sold", controllers.GetAllProjectsSold)
	r.GET("/projects/sold/:id", controllers.GetProjectsSoldByParamID)
	jwt.GET("/projects/buy", controllers.GetProjectsBuyByToken)

	r.GET("/projects/:id", controllers.GetProjectsById)
	r.GET("/projects/investee/:userid", controllers.GetProjectByUserId)
	jwt.GET("/investee/projects", controllers.GetProjectsByUserId)
	jwt.POST("/projects", controllers.CreateProjects)
	jwt.PATCH("/projects/:id", controllers.UpdateProjects)
	jwt.PATCH("/projects/status/:id", controllers.UpdateProjectsStatus)

	jwt.DELETE("/projects/:id", controllers.DeleteProjects)

	r.GET("/projects/some", controllers.GetSomeProjects)
	r.GET("/projects/filter", controllers.GetProjectsFiltered)

	jwt.GET("/projects/count", controllers.GetProjectsCount)

	//cartItems
	jwt.GET("/cartItems", controllers.GetAllCartItems)
	jwt.GET("/cartItems/:id", controllers.GetCartItemsById)
	jwt.POST("/cartItems", controllers.CreateCartItems)
	jwt.PATCH("/cartItems/:id", controllers.UpdateCartItems)
	jwt.DELETE("/cartItems/:id", controllers.DeleteCartItem)

	// payment
	jwt.POST("/payments", controllers.CreateMidtransTransaction)
	r.POST("/payments/notification", controllers.PaymentNotification)

	// Installment
	jwt.GET("/installment/:id", controllers.GetAllInstallmentByProjectID)
	jwt.PATCH("/installment/status/:id", controllers.UpdateInstallmentStatus)

	// swagger routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// routes not used
	// r.GET("/images", controllers.GetAllImages)
	// jwt.DELETE("/cart/:id", controllers.DeleteCart)

	return r
}
