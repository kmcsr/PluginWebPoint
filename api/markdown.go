
package api

import (
	"strings"
	"net/url"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

type Option struct{
	URLPrefix string
	HeadingIDPrefix string
	HeadingIDSuffix string
}

func (o *Option)exec(n ast.Node)(ast.Node){
	if len(o.URLPrefix) > 0 {
		if l, ok := n.(*ast.Link); ok {
			dest := (string)(l.Destination)
			if !strings.Contains(dest, "://") && !strings.HasPrefix(dest, "/") {
				out, err := url.JoinPath(o.URLPrefix, dest)
				if err != nil {
					out = "about:blank#err:" + err.Error()
				}
				l.Destination = ([]byte)(out)
				return l
			}
		}
	}
	return n
}

func (o *Option)Walk(n ast.Node)(ast.Node){
	n = o.exec(n)
	children := n.GetChildren()
	for i, c := range children {
		children[i] = o.Walk(c)
	}
	return n
}

func (o *Option)SetHtmlOptions(opt *html.RendererOptions){
	opt.HeadingIDPrefix = o.HeadingIDPrefix
	opt.HeadingIDSuffix = o.HeadingIDSuffix
}

func RenderMarkdown(src []byte, opt ...*Option)(out []byte){
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	opts := html.RendererOptions{
		Flags: html.CommonFlags | html.LazyLoadImages,
	}
	a := p.Parse(src)
	if len(opt) > 0 {
		opt[0].Walk(a)
		opt[0].SetHtmlOptions(&opts)
	}
	r := html.NewRenderer(opts)
	out1 := markdown.Render(a, r)
	out = bluemonday.UGCPolicy().SanitizeBytes(out1)
	return
}
