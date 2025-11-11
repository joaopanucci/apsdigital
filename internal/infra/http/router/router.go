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

func NewRouter(database *db.PostgresDB, cfg *config.Config) *gin.Engine {
	r := gin.New()

	// Add middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.CORSMiddleware())

	// Serve static files (uploaded PDFs)
	r.Static("/uploads", "./uploads")

	// Initialize repositories
	municipalityRepo := repositories.NewMunicipalityRepository(database)
	paymentRepo := repositories.NewPaymentRepository(database)
	resolutionRepo := repositories.NewResolutionRepository(database)

	// Initialize services
	// authService := services.NewAuthService(userRepo, refreshTokenRepo, roleRepo, cfg)  // TODO: Fix AuthorizeUser method
	municipalityService := services.NewMunicipalityService(municipalityRepo)
	// tabletService := services.NewTabletService(tabletRepo, userRepo)  // TODO: Fix AuthorizeUser method
	paymentService := services.NewPaymentService(paymentRepo)
	resolutionService := services.NewResolutionService(resolutionRepo)

	// Initialize controllers
	// authController := controllers.NewAuthController(authService)  // TODO: Fix AuthorizeUser method
	municipalityController := controllers.NewMunicipalityController(municipalityService)
	// tabletController := controllers.NewTabletController(tabletService)  // TODO: Fix AuthorizeUser method
	paymentController := controllers.NewPaymentController(paymentService)
	resolutionController := controllers.NewResolutionController(resolutionService)

	// API Routes with /api/v1 prefix (for direct API access)
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Municipalities
		municipalities := api.Group("/municipalities")
		{
			municipalities.GET("/", municipalityController.GetMunicipalities)
		}

		// Payments
		payments := api.Group("/payments")
		{
			payments.GET("/", paymentController.GetPaymentsPublic)
			payments.GET("/competences", paymentController.GetCompetencesPublic)
			payments.GET("/years", paymentController.GetYearsPublic)
			payments.POST("/upload", paymentController.UploadPaymentFile)
		}

		// Resolutions
		resolutions := api.Group("/resolutions")
		{
			resolutions.GET("/", resolutionController.GetResolutionsPublic)
			resolutions.GET("/types", resolutionController.GetTypesPublic)
			resolutions.GET("/years", resolutionController.GetYearsPublic)
			resolutions.POST("/upload", resolutionController.UploadResolutionFile)
		}
	}

	// Frontend API Routes (without /api/v1 prefix for frontend compatibility)
	// These will be accessed via nginx proxy as /api/ -> backend:8080/
	{
		// Health check
		r.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Municipalities
		municipalitiesRoot := r.Group("/municipalities")
		{
			municipalitiesRoot.GET("/", municipalityController.GetMunicipalities)
		}

		// Payments
		paymentsRoot := r.Group("/payments")
		{
			paymentsRoot.GET("/", paymentController.GetPaymentsPublic)
			paymentsRoot.GET("/competences", paymentController.GetCompetencesPublic)
			paymentsRoot.GET("/years", paymentController.GetYearsPublic)
			paymentsRoot.POST("/upload", paymentController.UploadPaymentFile)
		}

		// Resolutions
		resolutionsRoot := r.Group("/resolutions")
		{
			resolutionsRoot.GET("/", resolutionController.GetResolutionsPublic)
			resolutionsRoot.GET("/types", resolutionController.GetTypesPublic)
			resolutionsRoot.GET("/years", resolutionController.GetYearsPublic)
			resolutionsRoot.POST("/upload", resolutionController.UploadResolutionFile)
		}
	}

	// Temporary workaround: Add /api prefix routes directly in backend
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Municipalities
		apiMunicipalities := apiGroup.Group("/municipalities")
		{
			apiMunicipalities.GET("/", municipalityController.GetMunicipalities)
		}

		// Payments
		apiPayments := apiGroup.Group("/payments")
		{
			apiPayments.GET("/", paymentController.GetPaymentsPublic)
			apiPayments.GET("/competences", paymentController.GetCompetencesPublic)
			apiPayments.GET("/years", paymentController.GetYearsPublic)
			apiPayments.POST("/upload", paymentController.UploadPaymentFile)
		}

		// Resolutions
		apiResolutions := apiGroup.Group("/resolutions")
		{
			apiResolutions.GET("/", resolutionController.GetResolutionsPublic)
			apiResolutions.GET("/types", resolutionController.GetTypesPublic)
			apiResolutions.GET("/years", resolutionController.GetYearsPublic)
			apiResolutions.POST("/upload", resolutionController.UploadResolutionFile)
		}
	}

	return r
}
