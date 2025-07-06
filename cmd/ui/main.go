package main

import (
	"html/template"
	"net/http"

	"github.com/amrittb/choto-link-ui/internal/assets"
	"github.com/gin-gonic/gin"
)

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

	router.Run()
}
