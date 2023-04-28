package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Alert struct {
	Receiver string `form:"receiver" json:"receiver" binding:"required"`
	Status   string `form:"status" json:"status" binding:"required"`
}

type URLToken struct {
	ID string `uri:"id" binding:"required"`
}

func main() {
	// Load environment variables
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	// Set to release mode when in production environment
	if viper.Get("ENV_NAME") == "Production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Respond with simple pong
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "Pong")
	})

	// Redirect to projects GH page
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://github.com/hadret/forwardly-go")
	})

	// Forward POST from Alertmanager to GET in Uptime Kuma
	r.POST("/:id", func(c *gin.Context) {
		// Assign environment variables and structs
		var url_token URLToken
		var alert Alert
		kuma_tokens := viper.GetString("KUMA_TOKENS")
		kuma_url := viper.GetString("KUMA_URL")

		// Do the redirect or fail trying
		if c.ShouldBindUri(&url_token) == nil && c.ShouldBind(&alert) == nil {
			if strings.Contains(kuma_tokens, url_token.ID) {
				resp, err := http.Get(fmt.Sprintf("%v/%v", kuma_url, url_token.ID))
				if err != nil {
					fmt.Printf("Error when trying to reach Uptime Kuma: %v", err)
					return
				}
				c.String(resp.StatusCode, resp.Status)
			} else {
				c.String(401, "Unauthorized")
			}
		} else {
			c.String(404, "Not found")
		}
	})

	r.Run(":8000")
}
