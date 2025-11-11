package router

import (
	"github.com/gin-gonic/gin"
	"github.com/joaopanucci/apsdigital/internal/config"
	"github.com/joaopanucci/apsdigital/internal/domain/services"
	"github.com/joaopanucci/apsdigital/internal/infra/db"
	"github.com/joaopanucci/apsdigital/internal/infra/http/controllers"
	"github.com/joaopanucci/apsdigital/internal/infra/http/middlewares"
	"github.com/joaopanucci/apsdigital/internal/infra/repositories"
)

func SetupRoutes(database *db.PostgresDB, cfg *config.Config) *gin.Engine {
	r := gin.New()

	// Add middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.CORSMiddleware())

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database)
	roleRepo := repositories.NewRoleRepository(database)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(database)
	municipalityRepo := repositories.NewMunicipalityRepository(database)
	tabletRepo := repositories.NewTabletRepository(database)
	paymentRepo := repositories.NewPaymentRepository(database.DB)
	resolutionRepo := repositories.NewResolutionRepository(database.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo, refreshTokenRepo, cfg)
	municipalityService := services.NewMunicipalityService(municipalityRepo)
	tabletService := services.NewTabletService(tabletRepo, userRepo)
	paymentService := services.NewPaymentService(paymentRepo)
	resolutionService := services.NewResolutionService(resolutionRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	municipalityController := controllers.NewMunicipalityController(municipalityService)
	tabletController := controllers.NewTabletController(tabletService)
	paymentController := controllers.NewPaymentController(paymentService)
	resolutionController := controllers.NewResolutionController(resolutionService)

	// Public routes
	public := r.Group("/api/v1")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "message": "APSDigital API is running"})
		})

		// Auth routes
		auth := public.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)
			auth.POST("/logout", authController.Logout)
		}
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middlewares.AuthMiddleware(cfg))
	{
		// User profile
		protected.GET("/me", authController.Me)

		// Roles (Admin and Coordenador only)
		roles := protected.Group("/roles")
		roles.Use(middlewares.RequireMinLevel(2)) // Coordenador or higher
		{
			roles.GET("/", func(c *gin.Context) {
				// TODO: Implement role listing
				c.JSON(200, gin.H{"message": "Role listing - TODO"})
			})
		}

		// Users management
		users := protected.Group("/users")
		{
			// List users (with role-based filtering)
			users.GET("/", func(c *gin.Context) {
				// TODO: Implement user listing with filters based on role
				c.JSON(200, gin.H{"message": "User listing - TODO"})
			})

			// Get pending authorization (managers and above)
			users.GET("/pending", middlewares.RequireMinLevel(3), func(c *gin.Context) {
				// TODO: Implement pending users listing
				c.JSON(200, gin.H{"message": "Pending users - TODO"})
			})

			// Authorize user (based on hierarchy)
			users.POST("/:id/authorize", middlewares.RequireMinLevel(3), func(c *gin.Context) {
				// TODO: Implement user authorization
				c.JSON(200, gin.H{"message": "User authorization - TODO"})
			})

			// Reject user authorization
			users.POST("/:id/reject", middlewares.RequireMinLevel(3), func(c *gin.Context) {
				// TODO: Implement user rejection
				c.JSON(200, gin.H{"message": "User rejection - TODO"})
			})
		}

		// Professions (Admin and Coordenador only)
		professions := protected.Group("/professions")
		professions.Use(middlewares.RequireMinLevel(2))
		{
			professions.GET("/", func(c *gin.Context) {
				// TODO: Implement profession listing
				c.JSON(200, gin.H{"message": "Profession listing - TODO"})
			})
			professions.POST("/", func(c *gin.Context) {
				// TODO: Implement profession creation
				c.JSON(200, gin.H{"message": "Profession creation - TODO"})
			})
			professions.PUT("/:id", func(c *gin.Context) {
				// TODO: Implement profession update
				c.JSON(200, gin.H{"message": "Profession update - TODO"})
			})
			professions.DELETE("/:id", func(c *gin.Context) {
				// TODO: Implement profession deletion
				c.JSON(200, gin.H{"message": "Profession deletion - TODO"})
			})
		}

		// Tablet requests
		tablets := protected.Group("/tablets")
		{
			tablets.GET("/", func(c *gin.Context) {
				// TODO: Implement tablet request listing (filtered by role)
				c.JSON(200, gin.H{"message": "Tablet request listing - TODO"})
			})
			tablets.POST("/", func(c *gin.Context) {
				// TODO: Implement tablet request creation
				c.JSON(200, gin.H{"message": "Tablet request creation - TODO"})
			})
			tablets.PUT("/:id/approve", middlewares.RequireMinLevel(3), func(c *gin.Context) {
				// TODO: Implement tablet request approval
				c.JSON(200, gin.H{"message": "Tablet request approval - TODO"})
			})
			tablets.PUT("/:id/reject", middlewares.RequireMinLevel(3), func(c *gin.Context) {
				// TODO: Implement tablet request rejection
				c.JSON(200, gin.H{"message": "Tablet request rejection - TODO"})
			})
		}

		// Municipalities
		municipalities := protected.Group("/municipalities")
		{
			municipalities.GET("/", municipalityController.GetMunicipalities)
		}

		// Tablets
		tablets := protected.Group("/tablets")
		{
			tablets.GET("/search", tabletController.SearchAgent)
			tablets.POST("/request", tabletController.RequestTablet)
			tablets.GET("/requests", tabletController.GetTabletRequests)
			tablets.PUT("/requests/:id/approve", middlewares.RequireMinLevel(3), tabletController.ApproveRequest)
			tablets.PUT("/requests/:id/reject", middlewares.RequireMinLevel(3), tabletController.RejectRequest)
		}

		// Payments (Admin, Coordenador, and Gerente only)
		payments := protected.Group("/payments")
		payments.Use(middlewares.RequireMinLevel(3))
		{
			payments.GET("/", paymentController.GetPayments)
			payments.POST("/", paymentController.CreatePayment)
			payments.GET("/:id", paymentController.GetPaymentByID)
			payments.PUT("/:id", paymentController.UpdatePayment)
			payments.DELETE("/:id", paymentController.DeletePayment)
			payments.GET("/competences", paymentController.GetCompetences)
			payments.GET("/years", paymentController.GetYears)
			payments.GET("/:id/view", paymentController.ViewPDF)
			payments.GET("/:id/download", paymentController.DownloadPDF)
		}

		// Resolutions (Admin and Coordenador only)
		resolutions := protected.Group("/resolutions")
		resolutions.Use(middlewares.RequireMinLevel(2))
		{
			resolutions.GET("/", resolutionController.GetResolutions)
			resolutions.POST("/", resolutionController.CreateResolution)
			resolutions.GET("/:id", resolutionController.GetResolutionByID)
			resolutions.PUT("/:id", resolutionController.UpdateResolution)
			resolutions.DELETE("/:id", resolutionController.DeleteResolution)
			resolutions.GET("/types", resolutionController.GetTypes)
			resolutions.GET("/years", resolutionController.GetYears)
			resolutions.GET("/recent", resolutionController.GetRecentResolutions)
			resolutions.GET("/:id/view", resolutionController.ViewPDF)
			resolutions.GET("/:id/download", resolutionController.DownloadPDF)
		}

		// Professionals listing (all authenticated users can see)
		professionals := protected.Group("/professionals")
		{
			professionals.GET("/", func(c *gin.Context) {
				// TODO: Implement professional listing with filters
				c.JSON(200, gin.H{"message": "Professional listing - TODO"})
			})
		}
	}

	return r
}