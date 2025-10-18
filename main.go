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

	r.GET("/me", func(c *gin.Context) {
		fact, err := fetchCatFact()
		timestamp := time.Now().UTC().Format(time.RFC3339Nano)

		// Graceful fallback if external API fails
		if err != nil {
			log.Printf("Error fetching cat fact: %v", err)
			fact = "Cats are fascinating creatures — even when APIs fail."
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
	})

	log.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// fetchCatFact gets a random cat fact, with timeout and error handling
func fetchCatFact() (string, error) {
	client := resty.New().
		SetTimeout(5 * time.Second).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second)

	url := "https://catfact.ninja/fact"
	var factResp CatFactResponse

	resp, err := client.R().
		SetResult(&factResp).
		Get(url)
	if err != nil {
		return "", fmt.Errorf("network error: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return factResp.Fact, nil
}
