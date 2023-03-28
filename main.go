package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golobby/dotenv"
	"golang.org/x/exp/slices"
)

type Alert struct {
	Receiver string `form:"receiver" json:"receiver" binding:"required"`
	Status   string `form:"status" json:"status" binding:"required"`
}

type Config struct {
	Env_Name    string   `env:"ENV_NAME"`
	Kuma_Url    string   `env:"KUMA_URL"`
	Kuma_Tokens []string `env:"KUMA_TOKENS"`
}

type URLToken struct {
	ID string `uri:"id" binding:"required"`
}

func main() {
	// Load environment variables from file or fallback to global ones
	config := Config{}
	file, err := os.Open(".env")
	if err != nil {
		config.Env_Name = os.Getenv("ENV_NAME")
		config.Kuma_Url = os.Getenv("KUMA_URL")
		config.Kuma_Tokens = strings.Split(os.Getenv("KUMA_TOKENS"), ",")
	}
	err = dotenv.NewDecoder(file).Decode(&config)

	// Set to release mode when in production environment
	if config.Env_Name == "Production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Respond with simple pong
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
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

		// Do the redirect or fail trying
		if c.ShouldBindUri(&url_token) == nil && c.ShouldBind(&alert) == nil {
			if slices.Contains(config.Kuma_Tokens, url_token.ID) {
				resp, err := http.Get(fmt.Sprintf("%v/%v", config.Kuma_Url, url_token.ID))
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
