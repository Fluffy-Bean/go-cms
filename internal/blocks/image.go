package blocks

import (
	"bytes"
	"fmt"
	"html/template"
)

type ImageBlock struct {
	Image        string
	Alt          string
	AltAsCaption bool
}

func (b ImageBlock) Render() string {
	templ, err := template.ParseFiles("./templates/blocks/image.html")
	if err != nil {
		return fmt.Sprint(`<p class="compile-error">Error compiling image block</p>`)
	}

	var compiled bytes.Buffer
	err = templ.Execute(&compiled, map[string]any{
		"Image":        b.Image,
		"Alt":          b.Alt,
		"AltAsCaption": b.AltAsCaption,
	})
	if err != nil {
		return fmt.Sprint(`<p class="compile-error">Error compiling image block</p>`)
	}

	return compiled.String()
}
