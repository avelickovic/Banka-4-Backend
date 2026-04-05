package server

import (
	"context"
	stderrors "errors"
	"log"
	"net/http"
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/auth"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/handler"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/client"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/config"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/middleware"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/validator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/logging"
	_ "github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/docs"
)

func NewServer(lc fx.Lifecycle, cfg *config.Configuration, healthHandler *handler.HealthHandler, taxHandler *handler.TaxHandler, exchangeHandler *handler.ExchangeHandler, orderHandler *handler.OrderHandler, verifier auth.TokenVerifier, permProvider auth.PermissionProvider, userClient client.UserServiceClient) {
	r := gin.New()

	InitRouter(r, cfg)

	SetupRoutes(r, healthHandler, taxHandler, exchangeHandler, orderHandler, verifier, permProvider, userClient)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	RegisterServerLifecycle(lc, server)
}

func InitRouter(r *gin.Engine, cfg *config.Configuration) {
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.URLs.FrontendBaseURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(logging.Logger())
	r.Use(errors.ErrorHandler())

	validator.RegisterValidators()
}

func SetupRoutes(r *gin.Engine, healthHandler *handler.HealthHandler, taxHandler *handler.TaxHandler, exchangeHandler *handler.ExchangeHandler, orderHandler *handler.OrderHandler, verifier auth.TokenVerifier, permProvider auth.PermissionProvider, userClient client.UserServiceClient) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/health", healthHandler.Health)

		exchanges := api.Group("/exchanges")
		{
			exchanges.GET("", exchangeHandler.GetAll)
			exchanges.PATCH("/:micCode/toggle", exchangeHandler.ToggleTradingEnabled)
		}

		orders := api.Group("/orders")
		orders.Use(auth.Middleware(verifier, permProvider))
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.PATCH("/:id/approve", middleware.RequireSupervisor(userClient), orderHandler.ApproveOrder)
			orders.PATCH("/:id/decline", middleware.RequireSupervisor(userClient), orderHandler.DeclineOrder)
			orders.PATCH("/:id/cancel", orderHandler.CancelOrder)
		}
		tax := api.Group("/tax")
		tax.Use(auth.Middleware(verifier, permProvider))
		{
			tax.GET("", middleware.RequireSupervisor(userClient), taxHandler.ListTaxUsers)
			tax.POST("/collect", middleware.RequireSupervisor(userClient), taxHandler.CollectTaxes)
		}
	}
}

func RegisterServerLifecycle(lc fx.Lifecycle, server *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && !stderrors.Is(err, http.ErrServerClosed) {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
