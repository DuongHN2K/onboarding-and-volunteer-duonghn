package feature

import (
	"net/http"

	_ "github.com/cesc1802/onboarding-and-volunteer-service/docs"
	"github.com/cesc1802/onboarding-and-volunteer-service/feature/authentication/storage"
	authStorage "github.com/cesc1802/onboarding-and-volunteer-service/feature/authentication/storage"
	authTransport "github.com/cesc1802/onboarding-and-volunteer-service/feature/authentication/transport"
	authUsecase "github.com/cesc1802/onboarding-and-volunteer-service/feature/authentication/usecase"
	"github.com/cesc1802/onboarding-and-volunteer-service/feature/middleware"
	userStorage "github.com/cesc1802/onboarding-and-volunteer-service/feature/user/storage"
	userTransport "github.com/cesc1802/onboarding-and-volunteer-service/feature/user/transport"
	userUsecase "github.com/cesc1802/onboarding-and-volunteer-service/feature/user/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	appliIdentityStorage "github.com/cesc1802/onboarding-and-volunteer-service/feature/user_identity/storage"
	appliIdentityTransport "github.com/cesc1802/onboarding-and-volunteer-service/feature/user_identity/transport"
	appliIdentityUsecase "github.com/cesc1802/onboarding-and-volunteer-service/feature/user_identity/usecase"

	"github.com/cesc1802/share-module/system"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @host localhost:8080
// @BasePath /api/v1
func RegisterHandlerV1(mono system.Service) {
	router := mono.Router()
	secretKey := storage.GetSecretKey()
	router.Use(cors.Default())
	// add swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"data": "success",
		})
	})
	v1 := router.Group("/api/v1")
	// Initialize repository
	authRepo := authStorage.NewAuthenticationRepository(mono.DB())
	userRepo := userStorage.NewAdminRepository(mono.DB())
	applicantRepo := userStorage.NewApplicantRepository(mono.DB())
	applicantRequestRepo := userStorage.NewApplicantRequestRepository(mono.DB())
	applicantIdentityRepo := appliIdentityStorage.NewUserIdentityRepository(mono.DB())

	// Initialize usecase
	authUseCase := authUsecase.NewUserUsecase(authRepo, secretKey)
	userUseCase := userUsecase.NewAdminUsecase(userRepo)
	applicantUseCase := userUsecase.NewApplicantUsecase(applicantRepo)
	applicantRequestUseCase := userUsecase.NewApplicantRequestUsecase(applicantRequestRepo)
	applicantIdenityUseCase := appliIdentityUsecase.NewUserIdentityUsecase(applicantIdentityRepo)

	// Initialize handler
	authHandler := authTransport.NewAuthenticationHandler(authUseCase)
	userHandler := userTransport.NewAuthenticationHandler(userUseCase)
	applicantHandler := userTransport.NewApplicantHandler(applicantUseCase)
	applicantRequestHandler := userTransport.NewRequestHandler(applicantRequestUseCase)
	applicantIdentityHandler := appliIdentityTransport.NewUserIdentityHandler(applicantIdenityUseCase)

	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)

		auth.POST("/register", authHandler.Register)
	}

	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware(secretKey))
	{
		admin.GET("/list-request", userHandler.GetListRequest)
		admin.GET("/request/:id", userHandler.GetRequestById)
		admin.GET("/list-pending-request", userHandler.GetListPendingRequest)
		admin.GET("/pending-request/:id", userHandler.GetPendingRequestById)
		admin.POST("/approve-request/:id", userHandler.ApproveRequest)
		admin.POST("/reject-request/:id", userHandler.RejectRequest)
		admin.POST("/add-reject-notes/:id", userHandler.AddRejectNotes)
		admin.DELETE("/delete-request/:id", userHandler.DeleteRequest)
	}

	applicant := v1.Group("/applicant")
	{
		applicant.POST("/", applicantHandler.CreateApplicant)
		applicant.PUT("/:id", applicantHandler.UpdateApplicant)
		applicant.DELETE("/:id", applicantHandler.CreateApplicant)
		applicant.GET("/:id", applicantHandler.FindApplicantByID)
	}

	appliRequest := v1.Group("/applicant-request")
	{
		appliRequest.POST("/", applicantRequestHandler.CreateRequest)
	}

	appliIdentity := v1.Group("applicant-identity")
	{
		appliIdentity.POST("/", applicantIdentityHandler.CreateUserIdentity)
		appliIdentity.GET("/:id", applicantIdentityHandler.FindUserIdentity)
		appliIdentity.PUT("/:id", applicantIdentityHandler.UpdateUserIdentity)
	}

}
