package blocks

import (
	"bytes"
	"fmt"
	"html/template"
)

type CodeBlock struct {
	Code string
}

func (b CodeBlock) Render() string {
	templ, err := template.ParseFiles("./templates/blocks/code.html")
	if err != nil {
		return fmt.Sprint(`<p class="compile-error">Error compiling code block</p>`)
	}

	if b.Code == "" {
		return fmt.Sprint(`<p class="compile-error">Error compiling code block</p>`)
	}

	var compiled bytes.Buffer
	err = templ.Execute(&compiled, map[string]any{
		"Code": b.Code,
	})
	if err != nil {
		return fmt.Sprint(`<p class="compile-error">Error compiling code block</p>`)
	}

	return compiled.String()
}
