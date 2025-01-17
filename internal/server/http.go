package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-advanced/docs"
	"github.com/go-nunu/nunu-layout-advanced/internal/handler"
	"github.com/go-nunu/nunu-layout-advanced/internal/pkg/middleware"
	"github.com/go-nunu/nunu-layout-advanced/internal/pkg/response"
	"github.com/go-nunu/nunu-layout-advanced/pkg/jwt"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewServerHTTP(
	logger *log.Logger,
	jwt *jwt.JWT,
	userHandler handler.UserHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// swagger doc
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	r.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)

	// No route group has permission
	noAuthRouter := r.Group("/")
	{

		noAuthRouter.GET("/", func(ctx *gin.Context) {
			logger.WithContext(ctx).Info("hello")
			response.HandleSuccess(ctx, map[string]interface{}{
				":)": "Thank you for using nunu!",
			})
		})

		noAuthRouter.POST("/register", userHandler.Register)
		noAuthRouter.POST("/login", userHandler.Login)
	}
	// Non-strict permission routing group
	noStrictAuthRouter := r.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
	{
		noStrictAuthRouter.GET("/user", userHandler.GetProfile)
	}

	// Strict permission routing group
	strictAuthRouter := r.Group("/").Use(middleware.StrictAuth(jwt, logger))
	{
		strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
	}

	return r
}
