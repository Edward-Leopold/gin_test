package main

import (
	"fmt"
	"log"
	"os"
    "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
    "net/http"
)

func main() {
    // Создаем новый экземпляр роутера
    r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")

	fmt.Println("API Key:", apiKey)

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

    // Запускаем сервер на порту 8080
    r.Run(":8080")
}