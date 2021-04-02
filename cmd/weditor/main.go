package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"

	httpapi "github.com/dolfly/weditor/api/http"
	"github.com/dolfly/weditor/web"
)

func main() {
	app := cli.NewApp()
	app.Name = "weditor"
	app.Usage = ""
	app.Version = httpapi.Version
	app.Flags = []cli.Flag{
		&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Value: 17310},
		&cli.StringFlag{Name: "addr", Aliases: []string{"a"}, Value: "127.0.0.1"},
		&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}, Value: false},
	}
	app.Action = action
	app.Run(os.Args)
}
func action(c *cli.Context) error {
	addr := c.String("addr")
	port := c.Int("port")
	debug := c.Bool("debug")
	r := gin.Default()
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	//r.StaticFS("/static", http.Dir("./web/dist/static"))
	r.StaticFS("/static", web.Static())
	r.GET("/favicon.ico", web.Favicon())
	r.GET("/", web.Index())
	r.GET("/widget", web.Widget())
	r.GET("/quit", httpapi.ActionQuit)
	api := r.Group("/api", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "x-requested-with")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	})
	{
		r.GET("/ws/v1/python", httpapi.ActionPython)
		api.GET("/v1/version", httpapi.ActionVersion)
		api.POST("/v1/connect", httpapi.ActionConnect)
		api.Any("/v2/devices/*rurl", httpapi.ActionDevice)
		api.Any("/v1/devices/*rurl", httpapi.ActionDevice)
		api.POST("/v1/widgets", httpapi.ActionWidget)
	}
	return r.Run(fmt.Sprintf("%s:%d", addr, port))
}
