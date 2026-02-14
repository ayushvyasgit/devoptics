package proxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceProxy struct {
	targetURL string
}

func NewServiceProxy(targetURL string) *ServiceProxy {
	return &ServiceProxy{targetURL: targetURL}
}

func (p *ServiceProxy) Forward(c *gin.Context) {
	// Build target URL
	targetURL := p.targetURL + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Read request body
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Create new request
	req, err := http.NewRequest(
		c.Request.Method,
		targetURL,
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Add user context headers (if available)
	if userID, exists := c.Get("user_id"); exists {
		req.Header.Set("X-User-ID", userID.(string))
	}
	if email, exists := c.Get("email"); exists {
		req.Header.Set("X-User-Email", email.(string))
	}

	// Execute request
	client := &http.Client{Timeout: 30}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Copy status code
	c.Status(resp.StatusCode)

	// Copy response body
	io.Copy(c.Writer, resp.Body)
}