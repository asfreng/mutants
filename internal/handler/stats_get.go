package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Handler) StatsGet(c *gin.Context) {
	fmt.Printf("Handle GET Stats\n")

	stats, err := p.LogicHandler.GetDnaCount(c)
	if err != nil {
		fmt.Printf("Error getting stats: %s\n", err.Error())
		c.JSON(err.GetCode(), err)
		return
	}

	fmt.Printf("Response DNA Stats: %+v\n", stats)

	c.JSON(http.StatusOK, stats)
}
