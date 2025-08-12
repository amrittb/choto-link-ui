package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/amrittb/choto-link-ui/internal/assets"
	"github.com/gin-gonic/gin"
)

type CreateChotoReq struct {
	LongUrl string `json:"longUrl" binding:"required"`
}

func sanitizeAndValidate(input string) (string, error) {
	longUrl := strings.TrimSpace(input)
	if longUrl == "" {
		return "", errors.New("lamo link is empty")
	}

	if !strings.HasPrefix(longUrl, "http://") && !strings.HasPrefix(longUrl, "https://") {
		if !strings.Contains(longUrl, ".") || strings.Contains(longUrl, " ") {
			return "", errors.New("lamo link is invalid")
		}
		longUrl = "http://" + longUrl
	}

	parsedUrl, err := url.Parse(longUrl)
	if err != nil {
		return "", errors.New("lamo link is invalid")
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return "", errors.New("lamo link must use http or https")
	}

	if parsedUrl.Host == "" {
		return "", errors.New("lamo link is invalid")
	}

	host := strings.ToLower(parsedUrl.Host)
	if host == "localhost" || host == "127.0.0.1" || host == "0.0.0.0" {
		return "", errors.New("lamo link is invalid")
	}

	if host == "choto.link" || strings.HasSuffix(host, ".choto.link") {
		return "", errors.New("cannot shorten choto.link URLs")
	}

	if !strings.Contains(parsedUrl.Host, ".") {
		return "", errors.New("lamo link is invalid")
	}

	parsedUrl.Fragment = ""
	normalizedUrl := parsedUrl.String()

	_, err = url.ParseRequestURI(normalizedUrl)
	if err != nil {
		return "", errors.New("lamo link is invalid")
	}

	return normalizedUrl, nil
}

func main() {
	router := gin.Default()
	templ := template.Must(template.New("").ParseFS(assets.TemplateFS, "web/templates/*.tmpl"))

	router.SetHTMLTemplate(templ)

	// Read backend base URL from environment
	apiBaseUrl := os.Getenv("API_BASE_URL")
	if apiBaseUrl == "" {
		apiBaseUrl = "http://localhost:8081"
	}
	apiBaseUrl = strings.TrimRight(apiBaseUrl, "/")

	// Read public Choto domain for building full URLs
	chotoBaseUrl := os.Getenv("CHOTO_BASE_URL")
	if chotoBaseUrl == "" {
		chotoBaseUrl = "http://localhost:8080"
	}
	chotoBaseUrl = strings.TrimRight(chotoBaseUrl, "/")

	// Handle Static Files
	router.GET("/static/*filepath", func(c *gin.Context) {
		file := c.Param("filepath") // e.g., /main.css
		c.FileFromFS("web/static"+file, http.FS(assets.StaticFS))
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.POST("/create", func(c *gin.Context) {
		var reqBody CreateChotoReq
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Request.",
			})
			return
		}

		sanitizedUrl, err := sanitizeAndValidate(reqBody.LongUrl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Call backend API to create short link
		backendURL := apiBaseUrl + "/api/v1/choto"
		payload := []byte("{\"longUrl\":\"" + sanitizedUrl + "\"}")
		req, err := http.NewRequest(http.MethodPost, backendURL, bytes.NewBuffer(payload))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reach backend"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read backend response"})
			return
		}

		// Try to parse backend response and return in desired format
		var backend map[string]any
		if err := json.Unmarshal(body, &backend); err != nil {
			// If parsing fails, just proxy raw
			c.Data(resp.StatusCode, "application/json", body)
			return
		}

		if v, ok := backend["shortUrl"].(string); ok && v != "" {
			c.JSON(resp.StatusCode, gin.H{
				"shortUrl": chotoBaseUrl + "/" + v,
			})
			return
		}

		// If the backend did not return shortUrl, proxy as-is
		c.Data(resp.StatusCode, "application/json", body)
	})

	router.GET("/:shortUrl", func(c *gin.Context) {
		shortUrl := c.Param("shortUrl")
		if shortUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid short url"})
			return
		}

		// Call backend to resolve the long URL
		backendURL := apiBaseUrl + "/api/v1/choto/" + url.PathEscape(shortUrl)
		req, err := http.NewRequest(http.MethodGet, backendURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare request"})
			return
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reach backend"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read backend response"})
			return
		}

		if resp.StatusCode != http.StatusOK {
			// Proxy backend error status and body
			c.Data(resp.StatusCode, "application/json", body)
			return
		}

		var backend map[string]any
		if err := json.Unmarshal(body, &backend); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "invalid backend response"})
			return
		}

		longUrl, ok := backend["longUrl"].(string)
		if !ok || strings.TrimSpace(longUrl) == "" {
			c.JSON(http.StatusBadGateway, gin.H{"error": "backend did not return longUrl"})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, longUrl)
	})

	router.Run()
}
