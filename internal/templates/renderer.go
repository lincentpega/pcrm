package templates

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type Renderer struct {
	pageTemplates     map[string]*template.Template
	fragmentTemplates map[string]*template.Template
}

func NewRenderer() *Renderer {
	layoutPath := "internal/templates/layout.html"
	pageTemplates := parsePageTemplates(layoutPath)
	fragmentTemplates := parseFragmentTemplates()

	return &Renderer{
		pageTemplates:     pageTemplates,
		fragmentTemplates: fragmentTemplates,
	}
}

func parsePageTemplates(layoutPath string) map[string]*template.Template {
	pageTemplates := make(map[string]*template.Template)
	pagesPattern := "internal/templates/pages/*.html"

	pageFiles, err := filepath.Glob(pagesPattern)
	if err != nil {
		log.Fatalf("Failed to find page templates: %v", err)
	}

	for _, pageFile := range pageFiles {
		templateName := getTemplateName(pageFile)

		tmpl, err := template.ParseFiles(layoutPath, pageFile)
		if err != nil {
			log.Fatalf("Failed to parse page template %s: %v", templateName, err)
		}

		pageTemplates[templateName] = tmpl
	}

	return pageTemplates
}

func parseFragmentTemplates() map[string]*template.Template {
	fragmentTemplates := make(map[string]*template.Template)
	fragmentsPattern := "internal/templates/fragments/*.html"

	fragmentFiles, err := filepath.Glob(fragmentsPattern)
	if err != nil {
		log.Fatalf("Failed to find fragment templates: %v", err)
	}

	for _, fragmentFile := range fragmentFiles {
		templateName := getTemplateName(fragmentFile)

		tmpl, err := template.ParseFiles(fragmentFile)
		if err != nil {
			log.Fatalf("Failed to parse fragment template %s: %v", templateName, err)
		}

		fragmentTemplates[templateName] = tmpl
	}

	return fragmentTemplates
}

func getTemplateName(filePath string) string {
	fileName := filepath.Base(filePath)
	return strings.TrimSuffix(fileName, ".html")
}

func (r *Renderer) RenderPage(w http.ResponseWriter, templateName string, data interface{}) error {
	tmpl, exists := r.pageTemplates[templateName]
	if !exists {
		log.Printf("Page template '%s' not found", templateName)
		http.Error(w, "Page template not found", http.StatusInternalServerError)
		return nil
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Page template rendering error for '%s': %v, data: %+v", templateName, err, data)
		http.Error(w, "Failed to render page template", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (r *Renderer) RenderFragment(w http.ResponseWriter, templateName string, data interface{}) error {
	tmpl, exists := r.fragmentTemplates[templateName]
	if !exists {
		log.Printf("Fragment template '%s' not found", templateName)
		http.Error(w, "Fragment template not found", http.StatusInternalServerError)
		return nil
	}

	if err := tmpl.ExecuteTemplate(w, "content", data); err != nil {
		log.Printf("Fragment template rendering error for '%s': %v, data: %+v", templateName, err, data)
		http.Error(w, "Failed to render fragment template", http.StatusInternalServerError)
		return err
	}
	return nil
}
