package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupRouter creates a new Gin instance with the same route setup
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/me", func(c *gin.Context) {
		fact, err := fetchCatFact()
		timestamp := time.Now().UTC().Format(time.RFC3339Nano)

		if err != nil {
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
	return r
}

func TestGetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "success", response["status"])

	user, ok := response["user"].(map[string]interface{})
	assert.True(t, ok, "user field should be an object")
	assert.Equal(t, "eifenaike@gmail.com", user["email"])
	assert.Equal(t, "Ifenaike Elias Ayooluwa", user["name"])
	assert.Equal(t, "Go/Gin", user["stack"])

	assert.NotEmpty(t, response["fact"])

	timestamp, ok := response["timestamp"].(string)
	assert.True(t, ok, "timestamp should be a string")
	assert.True(t, isISO8601(timestamp), "timestamp must be valid ISO 8601 UTC format")
}

func isISO8601(s string) bool {
	_, err := time.Parse(time.RFC3339Nano, s)
	if err == nil {
		return true
	}
	isoRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z$`)
	return isoRegex.MatchString(s)
}
