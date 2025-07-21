package main

import (
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"strings"

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

		sanitizedUrl, err := sanitizeAndValidate(json.LongUrl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"shortUrl": sanitizedUrl,
		})
	})

	router.Run()
}
