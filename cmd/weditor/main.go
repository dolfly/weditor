package main

import (
	httpapi "github.com/dolfly/weditor/api/http"
	"github.com/dolfly/weditor/web"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//r.StaticFS("/static", http.Dir("./web/dist/static"))
	r.StaticFS("/static", web.Static())
	r.GET("/favicon.ico", web.Favicon())
	r.GET("/", web.Index())
	r.GET("/widget", web.Widget())
	r.GET("/quit", httpapi.Quit)
	api := r.Group("/api", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "x-requested-with")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	})
	{
		r.GET("/ws/v1/python", httpapi.Python)
		api.GET("/v1/version", httpapi.Version)
		api.POST("/v1/connect", httpapi.Connect)
		api.Any("/v2/devices/*rurl", httpapi.Device)
		api.Any("/v1/devices/*rurl", httpapi.Device)
		api.POST("/v1/widgets", httpapi.Widget)
	}
	r.Run(":8833")
}
