package main

import (
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/amrittb/choto-link-ui/internal/assets"
	"github.com/gin-gonic/gin"
)

type CreateChotoReq struct {
	LongUrl string `json:"longUrl" binding:"required"`
}

// isValidDomain checks if the URL has at least a second-level domain and TLD
func isValidDomain(urlStr string) bool {
	// Remove scheme if present for domain validation
	cleanUrl := strings.TrimPrefix(strings.TrimPrefix(urlStr, "http://"), "https://")
	
	// Remove path, query, and fragment
	if idx := strings.Index(cleanUrl, "/"); idx != -1 {
		cleanUrl = cleanUrl[:idx]
	}
	if idx := strings.Index(cleanUrl, "?"); idx != -1 {
		cleanUrl = cleanUrl[:idx]
	}
	if idx := strings.Index(cleanUrl, "#"); idx != -1 {
		cleanUrl = cleanUrl[:idx]
	}
	
	// Remove port if present
	if idx := strings.LastIndex(cleanUrl, ":"); idx != -1 {
		if portPart := cleanUrl[idx+1:]; regexp.MustCompile(`^\d+$`).MatchString(portPart) {
			cleanUrl = cleanUrl[:idx]
		}
	}
	
	// Check for valid domain format: at least one dot and valid characters
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.[a-zA-Z]{2,}(\.[a-zA-Z]{2,})*$`)
	return domainRegex.MatchString(cleanUrl)
}

func main() {
	router := gin.Default()
	templ := template.Must(template.New("").ParseFS(assets.TemplateFS, "web/templates/*.tmpl"))

	router.SetHTMLTemplate(templ)

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
		var json CreateChotoReq
		err := c.ShouldBindJSON(&json)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Request.",
			})
			return
		}

		longUrl := strings.Trim(json.LongUrl, " ")
		if longUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Empty Long Url.",
			})
			return
		}

		_, err = url.ParseRequestURI(longUrl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Long Url.",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"shortUrl": longUrl,
		})
	})

	router.Run()
}
