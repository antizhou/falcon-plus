package http

import (
	"fmt"
	"net/http"

	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/g"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/strategy"
	"github.com/open-falcon/falcon-plus/modules/agent/log_collector/worker"

	"github.com/gin-gonic/gin"
)

// Start http api
func Start() {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
	router.GET("/strategy", func(c *gin.Context) {
		c.JSON(http.StatusOK, strategy.GetListAll())
	})

	router.GET("/cached", func(c *gin.Context) {
		c.String(http.StatusOK, worker.GetCachedAll())
	})

	router.POST("/check", func(c *gin.Context) {
		log := c.PostForm("log")
		c.JSON(http.StatusOK, CheckLogByStrategy(log))
	})

	router.Run(fmt.Sprintf("0.0.0.0:%d", g.Conf().Http.HTTPPort))
}
