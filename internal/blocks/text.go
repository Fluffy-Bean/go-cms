package blocks

import (
	"bytes"
	"fmt"
	"html/template"
)

type TextBlock struct {
	Text string
}

func (b TextBlock) Render() string {
	templ, err := template.ParseFiles("./templates/blocks/text.html")
	if err != nil {
		return fmt.Sprint(`<p class="compile-error">Error compiling text block</p>`)
	}

	var compiled bytes.Buffer
	err = templ.Execute(&compiled, map[string]any{
		"Text": b.Text,
	})
	if err != nil {
		return fmt.Sprint(`<p class="compile-error">Error compiling text block</p>`)
	}

	return compiled.String()
}
