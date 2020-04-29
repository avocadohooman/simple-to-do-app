package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/auth0-community/go-auth0"
	"github.com/avocadohooman/handlers"
	"github.com/gin-gonic/gin"
	jose "gopkg.in/square/go-jose.v2"
)

var (
	audience     string
	domain       string
	ACCESS_TOKEN = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImFDRUFyVTByVHg4SHBvRUpJQTdtbiJ9.eyJpc3MiOiJodHRwczovL2Rldi1mYmo2dmRvZS5ldS5hdXRoMC5jb20vIiwic3ViIjoiUUxtVE1kbHhEcEVidWV2ZTNpRVQ5MzQ1elZmRnBGR2JAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vbXktZ29sYW5nLWFwaSIsImlhdCI6MTU4ODE2NzkwMywiZXhwIjoxNTg4MjU0MzAzLCJhenAiOiJRTG1UTWRseERwRWJ1ZXZlM2lFVDkzNDV6VmZGcEZHYiIsImd0eSI6ImNsaWVudC1jcmVkZW50aWFscyJ9.Ga0dYvsHwDwkBsfQjoXVsuAUMWxkjXy4Vtf7_shaIkS2vtVXfYWU4bXBzTI6pBRgU6UPrSTb7PLip1H8X1KJviHwhaPP0b5der-6oeuNv75JlkyXYq9O2AB94a_VWXpHilllLQm1xYeCZW7ow9YKn5LqXua6xe8Ujj4rRzWJGcm9e8h1DsC0LZ7E_vT26Lev5tGbh6SrRXRMVPpJ9VZcot4UoPR5OSgFTX8O1hLcC0poSwYwm6cvM0dcsNcvddunU5Vh4y_YcrKLnXfCmVz1QyGFoCq3ZahFNem2-M7Ip-GkCDbQwviY_rCk1aGtXAdY8WaZDNpqNkk5em2ei2KbKA"
)

func main() {
	setAuth0Variables()
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui" + path.Join(dir, file))
		}
	})

	authorized := r.Group("/")
	authorized.Use(authRequired())
	authorized.GET("/todo", handlers.GetTodoListHandler)
	authorized.POST("/todo", handlers.AddToHandler)
	authorized.DELETE("/todo/:id", handlers.DeleteTodoHandler)
	authorized.PUT("/todo", handlers.CompleteTodoHandler)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
}

func setAuth0Variables() {
	audience = os.Getenv("AUTH0_API_IDENTIFIER")
	domain = os.Getenv("AUTH0_DOMAIN")
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth0Domain = "https://" + domain + "/"
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: auth0Domain + ".well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{audience}, auth0Domain, jose.RS256)
		validator := auth0.NewValidator(configuration, nil)
		_, err := validator.ValidateRequest(c.Request)

		if err != nil {
			log.Println(err)
			terminateWithError(http.StatusUnauthorized, "token is not valid", c)
			return
		}
		c.Next()
	}
}

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
