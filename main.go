package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Alert struct {
	Receiver string `json:"receiver" binding:"required"`
	Status   string `json:"status" binding:"required"`
}

type URLToken struct {
	ID string `uri:"id" binding:"required"`
}

func main() {
	// Load environment variables
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	// Assign environment structs and variables
	var url_token URLToken
	var alert Alert
	admin_pass := viper.GetString("ADMIN_PASS")
	admin_user := viper.GetString("ADMIN_USER")
	env_name := viper.GetString("ENV_NAME")
	kuma_tokens := viper.GetString("KUMA_TOKENS")
	kuma_url := viper.GetString("KUMA_URL")

	// Set to release mode when in production environment
	if env_name == "Production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Start Gin with logging and recovery middleware
	r := gin.Default()

	// Respond with simple pong
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Pong")
	})

	// Admin endpoint returning set environment variables
	if admin_user != "" {
		r.GET("/admin", gin.BasicAuth(gin.Accounts{admin_user: admin_pass}), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"env_name":    env_name,
				"kuma_tokens": kuma_tokens,
				"kuma_url":    kuma_url,
			})
		})
	}

	// Redirect to projects GH page
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://github.com/hadret/forwardly-go")
	})

	// Forward POST from Alertmanager to GET in Uptime Kuma
	r.POST("/:id", func(c *gin.Context) {
		if c.ShouldBindUri(&url_token) == nil && c.ShouldBindJSON(&alert) == nil {
			if strings.Contains(kuma_tokens, url_token.ID) {
				resp, err := http.Get(fmt.Sprintf("%v/%v", kuma_url, url_token.ID))
				if err != nil {
					fmt.Printf("Error when trying to reach Uptime Kuma: %v", err)
					return
				}
				c.JSON(resp.StatusCode, resp.Status)
			} else {
				c.JSON(http.StatusUnauthorized, "Unauthorized")
			}
		} else {
			c.JSON(http.StatusBadRequest, "Not found")
		}
	})

	// Start the server
	r.Run(":8000")
}
