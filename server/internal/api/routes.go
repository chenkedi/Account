package api

import (
	"account/internal/api/handlers"
	"account/internal/api/middleware"
	"account/internal/business/services"
	"account/internal/data/repository"
	"account/internal/sync"
	"account/pkg/auth"
	"account/pkg/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func SetupRouter(cfg *config.Config, db *sqlx.DB, redis *redis.Client, logger *zap.Logger) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logging(logger))

	// Initialize repositories
	tokenMgr := auth.NewTokenManager(cfg)
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	syncRepo := repository.NewSyncRepository(db, accountRepo, categoryRepo, transactionRepo)

	// Initialize sync engine
	syncEngine := sync.NewSyncEngine(syncRepo, accountRepo, categoryRepo, transactionRepo, logger)

	// Initialize services
	accountService := services.NewAccountService(accountRepo, categoryRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo, accountRepo, categoryRepo)
	syncService := services.NewSyncService(syncRepo, accountRepo, categoryRepo, transactionRepo)
	importService := services.NewImportService(transactionRepo, accountRepo, categoryRepo, logger)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, categoryRepo, tokenMgr, logger)
	accountHandler := handlers.NewAccountHandler(accountService, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger)
	transactionHandler := handlers.NewTransactionHandler(transactionService, logger)
	syncHandler := handlers.NewSyncHandler(syncService, logger)
	importHandler := handlers.NewImportHandler(importService, logger)
	wsHandler := handlers.NewWebSocketHandler(syncEngine.GetNotifier(), tokenMgr, logger)

	authMiddleware := middleware.NewAuthMiddleware(tokenMgr)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes (no auth required)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// Import info (no auth required for template info)
		importInfo := api.Group("/import")
		{
			importInfo.GET("/sources", importHandler.GetSupportedSources)
			importInfo.GET("/template", importHandler.GetTemplateInfo)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			// User endpoints
			protected.GET("/me", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				email, _ := c.Get("email")
				c.JSON(http.StatusOK, gin.H{
					"user_id": userID,
					"email":   email,
				})
			})

			// Account endpoints
			accounts := protected.Group("/accounts")
			{
				accounts.POST("", accountHandler.CreateAccount)
				accounts.GET("", accountHandler.GetAllAccounts)
				accounts.GET("/:id", accountHandler.GetAccount)
				accounts.PUT("/:id", accountHandler.UpdateAccount)
				accounts.DELETE("/:id", accountHandler.DeleteAccount)
			}

			// Category endpoints
			categories := protected.Group("/categories")
			{
				categories.POST("", categoryHandler.CreateCategory)
				categories.GET("", categoryHandler.GetAllCategories)
				categories.GET("/type/:type", categoryHandler.GetCategoriesByType)
				categories.GET("/:id", categoryHandler.GetCategory)
				categories.PUT("/:id", categoryHandler.UpdateCategory)
				categories.DELETE("/:id", categoryHandler.DeleteCategory)
			}

			// Transaction endpoints
			transactions := protected.Group("/transactions")
			{
				transactions.POST("", transactionHandler.CreateTransaction)
				transactions.GET("", transactionHandler.GetAllTransactions)
				transactions.GET("/range", transactionHandler.GetTransactionsByDateRange)
				transactions.GET("/stats", transactionHandler.GetStats)
				transactions.GET("/stats/detailed", transactionHandler.GetDetailedStats)
				transactions.GET("/:id", transactionHandler.GetTransaction)
				transactions.PUT("/:id", transactionHandler.UpdateTransaction)
				transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
			}

			// Sync endpoints
			sync := protected.Group("/sync")
			{
				sync.POST("/pull", syncHandler.Pull)
				sync.POST("/push", syncHandler.Push)
			}

			// Import endpoints
			importGroup := protected.Group("/import")
			{
				importGroup.POST("/upload", importHandler.UploadAndParse)
				importGroup.POST("/execute", importHandler.ExecuteImport)
			}
		}
	}

	// WebSocket endpoint (doesn't use standard auth middleware, token in query)
	router.GET("/ws/sync", wsHandler.HandleWebSocket)

	return router
}
