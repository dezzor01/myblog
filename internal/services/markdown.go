package services

import (
	"html/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func RenderMarkdown(content string) template.HTML {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Tables
	p := parser.NewWithExtensions(extensions)

	renderer := html.NewRenderer(html.RendererOptions{
		Flags: html.CommonFlags | html.HrefTargetBlank | html.TOC,
	})

	return template.HTML(markdown.ToHTML([]byte(content), p, renderer))
}
