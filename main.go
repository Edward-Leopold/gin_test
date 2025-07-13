package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Создаем новый экземпляр роутера
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("FATAL: API_KEY environment variable not set. Check your .env file and docker-compose.yaml configuration.")
	}

	fmt.Println("Successfully loaded API Key.")

	// Определяем маршрут для главной страницы
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Привет, Gin!")
	})

	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name") // Получаем параметр name из URL
		c.String(http.StatusOK, "Привет, %s!", name)
	})

	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	r.POST("/user", func(c *gin.Context) {
		var user User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Пользователь получен",
			"user":    user,
		})
	})

	type Message struct {
		Message string `json"message"`
	}

	r.POST("/api/message", func(c *gin.Context) {
		var msg Message

		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//{"chat":      {"model":"gpt-3.5-turbo","messages":[{"role":"user","content":"hello"}]}        }
		postBody, _ := json.Marshal(map[string]interface{}{
			"chat": map[string]interface{}{
				"model": "gpt-3.5-turbo",
				"messages": []map[string]string{
					{
						"role":    "user",
						"content": msg.Message,
					},
				},
			},
		})
		requestBody := bytes.NewBuffer(postBody)

		req, err := http.NewRequest(
			"POST",
			"https://caila.io/api/mlpgate/account/just-ai/model/openai-proxy/predict",
			requestBody,
		)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		req.Header.Set("MLP-API-KEY", apiKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()
		//Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		sb := string(body)
		log.Printf(sb)

		c.JSON(http.StatusOK, sb)
	})

	// Запускаем сервер на порту 8080
	r.Run(":8080")
}
