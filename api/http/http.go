package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dolfly/uiautomator"
	"github.com/dolfly/weditor/api/util"
	"github.com/gin-gonic/gin"
)

func Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":    "weditor",
		"version": "0.0.6",
	})
}

func Connect(c *gin.Context) {
	platform := strings.ToLower(c.PostForm("platform"))
	deviceurl := c.PostForm("deviceUrl")
	_ = func(platform, url string) string {
		return ""
	}(platform, deviceurl)
	ua := uiautomator.Connect(deviceurl)
	status, err := ua.Ping()
	if err != nil {
		fmt.Println(status)
	}
	screenWebSocketUrl := ""
	if platform == "android" {
		if u, err := url.Parse(deviceurl); err == nil {
			switch strings.ToLower(u.Scheme) {
			case "http":
				u.Scheme = "ws"
			case "https":
				u.Scheme = "wss"
			}
			u.Path = "/minicap"
			screenWebSocketUrl = u.String()
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"deviceId":           deviceurl,
		"platform":           platform,
		"success":            status == "pong",
		"status":             status,
		"screenWebSocketUrl": screenWebSocketUrl,
	})
}

func Device(c *gin.Context) {
	rurl := c.Param("rurl")
	remote, err := url.Parse(strings.TrimLeft(rurl, "/"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "error",
			"errmsg": err.Error(),
		})
		return
	}
	durl := fmt.Sprintf("%s://%s", remote.Scheme, remote.Host)
	switch strings.TrimLeft(strings.ToLower(remote.Path), "/") {
	case "hierarchy":
		c.JSON(http.StatusOK, hierarchy(durl))
		return
	case "screenshot":
		c.JSON(http.StatusOK, screenshot(durl))
		return
	case "exec":
		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"duration": 22222,
			"context":  "content",
		})
		return
	default:
	}
}
func Widget(c *gin.Context) {
	req := map[string]interface{}{}
	c.BindJSON(&req)
	c.JSON(http.StatusOK, req)
}
func Quit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
func screenshot(durl string) gin.H {
	// ua := uiautomator.Connect(durl)
	// screenshot, _ := ua.GetScreenshot()
	res, err := http.Get(durl + "/screenshot/0")
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	obj := gin.H{
		"type":     "jpeg",
		"encoding": "base64",
		"data":     base64.StdEncoding.EncodeToString(data),
	}
	return obj
}

func hierarchy(durl string) map[string]interface{} {
	ua := uiautomator.Connect(durl)
	ai, err := ua.GetCurrentApp()
	if err != nil {
		return nil
	}
	ws, err := ua.GetWindowSize()
	if err != nil {
		return nil
	}

	res, err := http.Get(durl + "/dump/hierarchy")
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	obj := struct {
		Id      int    `json:"id"`
		JsonRpc string `json:"jsonrpc"`
		Result  string `json:"result"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&obj)
	if err != nil {
		return nil
	}
	x, err := util.Convert([]byte(obj.Result))
	if err != nil {
		return nil
	}
	return map[string]interface{}{
		"xmlHierarchy":  obj.Result,
		"jsonHierarchy": x.Hierarchy,
		"activity":      ai.Activity,
		"packageName":   ai.Package,
		"windowSize":    []int{ws.Width, ws.Height},
	}
}
