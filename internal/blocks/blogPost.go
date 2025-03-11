package blocks

import (
	"bytes"
	"fmt"
	"html/template"
)

type BlogPostBlock struct {
	Title       string
	Summary     string
	PublishDate string
}

func (b BlogPostBlock) Render() string {
	templ, err := template.ParseFiles("./templates/blocks/blog_post.html")
	if err != nil {
		return fmt.Sprint(`<p class="compile-error">Error compiling blogPost block</p>`)
	}

	var compiled bytes.Buffer
	err = templ.Execute(&compiled, map[string]any{
		"Title":       b.Title,
		"Summary":     b.Summary,
		"PublishDate": b.PublishDate,
	})
	if err != nil {
		return fmt.Sprint(`<p class="compile-error">Error compiling blogPost block</p>`)
	}

	return compiled.String()
}
