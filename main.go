package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type CatFactResponse struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Stack string `json:"stack"`
}

type APIResponse struct {
	Status    string `json:"status"`
	User      User   `json:"user"`
	Timestamp string `json:"timestamp"`
	Fact      string `json:"fact"`
}

func main() {
	r := gin.Default()
	r.GET("/me", getProfile)
	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}

func getProfile(c *gin.Context) {
	fact, err := fetchCatFact()
	timestamp := time.Now().UTC().Format(time.RFC3339Nano)

	if err != nil {
		log.Printf("Error fetching cat fact: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"status":    "error",
			"message":   "Failed to fetch cat fact from external API",
			"timestamp": timestamp,
		})
		return
	}

	response := APIResponse{
		Status: "success",
		User: User{
			Email: "eifenaike@gmail.com",
			Name:  "Ifenaike Elias Ayooluwa",
			Stack: "Go/Gin",
		},
		Timestamp: timestamp,
		Fact:      fact,
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
}

func fetchCatFact() (string, error) {
	client := resty.New().
		SetTimeout(5 * time.Second).
		SetRetryCount(1).
		SetRetryWaitTime(500 * time.Millisecond)

	url := "https://catfact.ninja/fact"
	var factResp CatFactResponse

	resp, err := client.R().SetResult(&factResp).Get(url)
	if err != nil {
		return "", fmt.Errorf("network error: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return factResp.Fact, nil
}
