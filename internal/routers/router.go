package routers

import (
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger" //nolint: goimports
	"github.com/swaggo/gin-swagger/swaggerFiles"

	// import swagger handler
	_ "github.com/go-eagle/eagle/api/http" // docs is generated by Swag CLI, you have to import it.

	"github.com/go-eagle/eagle/internal/handler/v1/user"
	mw "github.com/go-eagle/eagle/internal/middleware"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/middleware"
)

// NewRouter loads the middlewares, routes, handlers.
func NewRouter() *gin.Engine {
	g := gin.New()
	// Use middleware
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middleware.Logging())
	g.Use(middleware.RequestID())
	g.Use(middleware.Metrics(app.Conf.Name))
	g.Use(middleware.Tracing(app.Conf.Name))
	g.Use(middleware.Timeout(3 * time.Second))
	g.Use(mw.Translations())

	// load web router
	LoadWebRouter(g)

	// 404 Handler.
	g.NoRoute(app.RouteNotFound)
	g.NoMethod(app.RouteNotFound)

	// swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// pprof router Profiling Routing
	// It is closed by default and can be opened in the development environment
	// Access: HOST/debug/pprof
	// Generate profile via HOST/debug/pprof/profile
	// View profile graph go tool pprof -http=:5000 profile
	// see: https://github.com/gin-contrib/pprof
	if app.Conf.EnablePprof {
		pprof.Register(g)
	}

	// HealthCheck Health Check Routing
	g.GET("/health", app.HealthCheck)
	// metrics router can be monitored in prometheus
	// Visually view the monitoring data of prometheus through grafana, use plugin 6671 to view
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// v1 router
	apiV1 := g.Group("/v1")
	apiV1.Use()
	{
		// Authentication related routes
		apiV1.POST("/register", user.Register)
		apiV1.POST("/login", user.Login)
		apiV1.POST("/login/phone", user.PhoneLogin)
		apiV1.GET("/vcode", user.VCode)

		// user
		apiV1.GET("/users/:id", user.Get)
		apiV1.Use(middleware.Auth())
		{
			apiV1.PUT("/users/:id", user.Update)
			apiV1.POST("/users/follow", user.Follow)
			apiV1.POST("/users/unfollow", user.Unfollow)
			apiV1.GET("/users/:id/following", user.FollowList)
			apiV1.GET("/users/:id/followers", user.FollowerList)
		}
	}

	return g
}