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

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/me", getProfile)
	return r
}

func TestGetProfile(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "success", response["status"])

	user, ok := response["user"].(map[string]interface{})
	assert.True(t, ok, "user field should be an object")
	assert.NotEmpty(t, user["email"])
	assert.NotEmpty(t, user["name"])
	assert.NotEmpty(t, user["stack"])

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
	// accept slightly shorter (RFC3339) too
	isoRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z$`)
	return isoRegex.MatchString(s)
}
