package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"challenge/internal/api"
	"challenge/internal/handler"
	"challenge/internal/util"
)

// SetupRouter .
func SetupRouter(handler *handler.Handler, config util.Config) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	// create a new serve mux and register the handlers
	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("/mutant", handler.MutantsPost)
	router.GET("/stats", handler.StatsGet)

	router.NoRoute(notFoundHandler)
	return router
}

func notFoundHandler(c *gin.Context) {
	fmt.Println("Route not found")
	err := &api.ErrorNotFound{}
	c.JSON(err.GetCode(), err)
}
